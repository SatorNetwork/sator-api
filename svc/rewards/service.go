package rewards

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/SatorNetwork/sator-api/lib/db"
	"github.com/SatorNetwork/sator-api/svc/qrcodes"
	"github.com/SatorNetwork/sator-api/svc/rewards/consts"
	"github.com/SatorNetwork/sator-api/svc/rewards/repository"
	"github.com/SatorNetwork/sator-api/svc/rewards/worker"
	"github.com/SatorNetwork/sator-api/svc/wallet"
)

//go:generate mockgen -destination=mock_repository.go -package=rewards github.com/SatorNetwork/sator-api/svc/rewards RewardsRepository
type (
	// Service struct
	Service struct {
		repo              RewardsRepository
		ws                walletService
		getLocker         db.GetLocker
		assetName         string
		explorerURLTmpl   string
		holdRewardsPeriod time.Duration
		minAmountToClaim  float64 // minimum amount to claim rewards
		worker            inProgressWorker
	}

	Winner struct {
		UserID uuid.UUID
		Points int
	}

	RewardsRepository interface {
		AddTransaction(ctx context.Context, arg repository.AddTransactionParams) error
		Withdraw(ctx context.Context, uid uuid.UUID) error
		GetTotalAmount(ctx context.Context, userID uuid.UUID) (float64, error)
		GetTransactionsByUserIDPaginated(ctx context.Context, arg repository.GetTransactionsByUserIDPaginatedParams) ([]repository.Reward, error)
		// GetAmountAvailableToWithdraw(ctx context.Context, arg repository.GetAmountAvailableToWithdrawParams) (float64, error)
		GetScannedQRCodeByUserID(ctx context.Context, arg repository.GetScannedQRCodeByUserIDParams) (repository.Reward, error)
		RequestTransactionsByUserID(ctx context.Context, userID uuid.UUID) error
		SetInProgressTransaction(ctx context.Context, arg repository.SetInProgressTransactionParams) error
		UpdateTransactionStatusByTxHash(ctx context.Context, arg repository.UpdateTransactionStatusByTxHashParams) error
	}

	ClaimRewardsResult struct {
		DisplayAmount   string  `json:"amount"`
		TransactionURL  string  `json:"transaction_url"`
		Amount          float64 `json:"-"`
		TransactionHash string  `json:"-"`
	}

	walletService interface {
		WithdrawRewards(ctx context.Context, userID uuid.UUID, amount float64) (string, error)
	}

	inProgressWorker interface {
		AddTransaction(transaction worker.InProgressTransactionStatusWorkerJob)
		Start()
	}

	// Option func to set custom service options
	Option func(*Service)
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(
	repo RewardsRepository,
	ws walletService,
	getLocker db.GetLocker,
	worker *worker.InProgressTransactionStatusWorker,
	opt ...Option,
) *Service {
	s := &Service{
		repo:              repo,
		ws:                ws,
		getLocker:         getLocker,
		assetName:         "SAO",
		explorerURLTmpl:   "https://explorer.solana.com/tx/%s?cluster=devnet",
		holdRewardsPeriod: time.Hour * 24 * 30,
		minAmountToClaim:  0,
		worker:            worker,
	}

	for _, fn := range opt {
		fn(s)
	}

	return s
}

func (s *Service) GetRewardsWallet(ctx context.Context, userID, walletID uuid.UUID) (wallet.Wallet, error) {
	totalRewards, _, err := s.GetUserRewards(ctx, userID)
	if err != nil {
		return wallet.Wallet{}, fmt.Errorf("could  not get rewards wallet: %w", err)
	}

	return wallet.Wallet{
		ID:    walletID.String(),
		Order: 99,
		Balance: []wallet.Balance{
			{
				Currency: "UNCLAIMED",
				Amount:   totalRewards,
			},
			// {
			// 	Currency: "Available to claim",
			// 	Amount:   availableRewards,
			// },
		},
		Actions: []wallet.Action{{
			Type: wallet.ActionClaimRewards.String(),
			Name: wallet.ActionClaimRewards.Name(),
			URL:  "rewards/claim",
		}},
	}, nil
}

// AddTransaction ...
func (s *Service) AddTransaction(ctx context.Context, uid, relationID uuid.UUID, relationType string, amount float64, trType int32) error {
	if err := s.repo.AddTransaction(ctx, repository.AddTransactionParams{
		UserID:          uid,
		RelationID:      uuid.NullUUID{UUID: relationID, Valid: true},
		Amount:          amount,
		TransactionType: trType,
		RelationType:    sql.NullString{String: relationType, Valid: true},
	}); err != nil {
		return fmt.Errorf("could not add transaction for user_id=%s, relation_id=%s, relation_type=%s: %w", uid.String(), relationID.String(), relationType, err)
	}

	return nil
}

// ClaimRewards send rewards to user by it and sets them to withdrawn.
func (s *Service) ClaimRewards(ctx context.Context, uid uuid.UUID) (ClaimRewardsResult, error) {
	// id := fmt.Sprintf("claim-rewards-%v", uid.String())
	// lock, err := s.getLocker.GetLock(ctx, id)
	// if err != nil {
	// 	return ClaimRewardsResult{}, fmt.Errorf("can't get lock by id: %v, err: %v", id, err)
	// }

	// ok, err := lock.Lock(ctx)
	// if err != nil {
	// 	return ClaimRewardsResult{}, fmt.Errorf("can't acquire a lock with id: %v, err: %v", id, err)
	// }
	// if !ok {
	// 	return ClaimRewardsResult{}, fmt.Errorf("lock %v is already acquired", id)
	// }

	amount, err := s.repo.GetTotalAmount(ctx, uid)
	if err != nil {
		if db.IsNotFoundError(err) {
			return ClaimRewardsResult{}, ErrRewardsAlreadyClaimed
		}

		return ClaimRewardsResult{}, fmt.Errorf("could not get total amount of rewards: %w", err)
	}

	if amount < s.minAmountToClaim {
		return ClaimRewardsResult{}, fmt.Errorf("%w: %.2f", ErrNotEnoughBalance, s.minAmountToClaim)
	}

	err = s.repo.RequestTransactionsByUserID(ctx, uid)
	if err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("%w", ErrInternalServerError)
	}

	txHash, err := s.ws.WithdrawRewards(ctx, uid, amount)
	if err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("could not create blockchain transaction: %w", err)
	}

	if err := s.repo.AddTransaction(ctx, repository.AddTransactionParams{
		UserID:          uid,
		Amount:          amount,
		TransactionType: consts.TransactionTypeWithdraw,
		TxHash:          sql.NullString{String: txHash, Valid: true},
		Status:          consts.TransactionStatusInProgress.String(),
	}); err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("could not add reward: %w", err)
	}

	s.worker.AddTransaction(worker.InProgressTransactionStatusWorkerJob{
		UserID: uid,
		Amount: amount,
		TxHash: txHash,
	})

	// We should release a lock in any case, even if context was cancelled
	// TODO(evg): get timeout from config
	// ctxt, _ := context.WithTimeout(context.Background(), 15 * time.Second)
	// if err := lock.Unlock(ctxt); err != nil {
	// 	return ClaimRewardsResult{}, fmt.Errorf("can't release a lock with id: %v, err: %v", id, err)
	// }

	return ClaimRewardsResult{
		Amount:          amount,
		DisplayAmount:   fmt.Sprintf("%.2f %s", amount, s.assetName),
		TransactionHash: txHash,
		TransactionURL:  fmt.Sprintf(s.explorerURLTmpl, txHash),
	}, nil
}

// GetUserRewards returns users available balance.
func (s *Service) GetUserRewards(ctx context.Context, uid uuid.UUID) (total float64, available float64, err error) {
	total, err = s.repo.GetTotalAmount(ctx, uid)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return 0, 0, fmt.Errorf("could not get total amount of rewards: %w", err)
		}
	}

	// available, err = s.repo.GetAmountAvailableToWithdraw(ctx, repository.GetAmountAvailableToWithdrawParams{
	// 	UserID:       uid,
	// 	NotAfterDate: time.Now().Add(-s.holdRewardsPeriod),
	// })
	// if err != nil {
	// 	if !db.IsNotFoundError(err) {
	// 		return 0, 0, fmt.Errorf("could not get available amount of rewards: %w", err)
	// 	}
	// }

	return total, 0, nil
}

// GetTransactions returns list of transactions from rewards wallet.
func (s *Service) GetTransactions(ctx context.Context, userID, walletID uuid.UUID, limit, offset int32) (wallet.Transactions, error) {
	txList, err := s.repo.GetTransactionsByUserIDPaginated(ctx, repository.GetTransactionsByUserIDPaginatedParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset})
	if err != nil {
		if db.IsNotFoundError(err) {
			return wallet.Transactions{}, nil
		}
		return nil, fmt.Errorf("could not get rewards transactions list: %w", err)
	}

	result := wallet.Transactions{}
	for _, tx := range txList {
		amount := tx.Amount
		if tx.TransactionType == consts.TransactionTypeWithdraw {
			amount = amount * (-1)
		}
		desc := tx.RelationType.String
		if desc == "" {
			desc = "claim rewards"
		}
		result = append(result, wallet.Transaction{
			ID:        tx.ID.String(),
			WalletID:  walletID.String(),
			TxHash:    desc,
			Amount:    amount,
			CreatedAt: tx.CreatedAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

// IsQRCodeScanned returns true if user got reward by this qrcode_id.
func (s *Service) IsQRCodeScanned(ctx context.Context, userID, qrcodeID uuid.UUID) (bool, error) {
	_, err := s.repo.GetScannedQRCodeByUserID(ctx, repository.GetScannedQRCodeByUserIDParams{
		UserID: userID,
		RelationID: uuid.NullUUID{
			UUID:  qrcodeID,
			Valid: true,
		},
		RelationType: sql.NullString{
			String: qrcodes.RelationTypeQRcodes,
			Valid:  true,
		},
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return false, nil
		}

		return false, fmt.Errorf("could not get rewards transactions list: %w", err)
	}

	return true, nil
}

func (s *Service) StartWorker() {
	s.worker.Start()
}
