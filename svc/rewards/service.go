package rewards

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/svc/rewards/repository"

	"github.com/google/uuid"
)

const (
	// TransactionTypeDeposit indicates that transaction type deposit.
	TransactionTypeDeposit = iota + 1
	// TransactionTypeWithdraw indicates that transaction type withdraw.
	TransactionTypeWithdraw
)

type (
	// Service struct
	Service struct {
		repo            rewardsRepository
		ws              walletService
		assetName       string
		explorerURLTmpl string
	}

	Winner struct {
		UserID uuid.UUID
		Points int
	}

	rewardsRepository interface {
		AddTransaction(ctx context.Context, arg repository.AddTransactionParams) error
		Withdraw(ctx context.Context, userID uuid.UUID) error
		GetTotalAmount(ctx context.Context, userID uuid.UUID) (float64, error)
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

	// Option func to set custom service options
	Option func(*Service)
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo rewardsRepository, ws walletService, opt ...Option) *Service {
	s := &Service{
		repo:            repo,
		ws:              ws,
		assetName:       "SAO",
		explorerURLTmpl: "https://explorer.solana.com/tx/%s?cluster=devnet",
	}

	for _, fn := range opt {
		fn(s)
	}

	return s
}

// AddTransaction ..
func (s *Service) AddTransaction(ctx context.Context, uid uuid.UUID, amount float64, qid uuid.UUID, trType int32) error {
	if err := s.repo.AddTransaction(ctx, repository.AddTransactionParams{
		UserID:          uid,
		QuizID:          qid,
		Amount:          amount,
		TransactionType: trType,
	}); err != nil {
		return fmt.Errorf("could not add transaction for user_id=%s and quiz_id=%s: %w", uid.String(), qid.String(), err)
	}

	return nil
}

// ClaimRewards send rewards to user by it and sets them to withdrawn.
func (s *Service) ClaimRewards(ctx context.Context, uid uuid.UUID) (ClaimRewardsResult, error) {
	amount, err := s.repo.GetTotalAmount(ctx, uid)
	if err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("could not get total amount of rewards: %w", err)
	}

	txHash, err := s.ws.WithdrawRewards(ctx, uid, amount)
	if err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("could not create blockchain transaction: %w", err)
	}

	if err = s.repo.Withdraw(ctx, uid); err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("could not update rewards status: %w", err)
	}

	if err := s.repo.AddTransaction(ctx, repository.AddTransactionParams{
		UserID:          uid,
		Amount:          amount,
		TransactionType: TransactionTypeWithdraw,
	}); err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("could not add reward: %w", err)
	}

	return ClaimRewardsResult{
		Amount:          amount,
		DisplayAmount:   fmt.Sprintf("%.2f %s", amount, s.assetName),
		TransactionHash: txHash,
		TransactionURL:  fmt.Sprintf(s.explorerURLTmpl, txHash),
	}, nil
}

// GetUserRewards returns users available balance.
func (s *Service) GetUserRewards(ctx context.Context, uid uuid.UUID) (float64, error) {
	amount, err := s.repo.GetTotalAmount(ctx, uid)
	if err != nil {
		return 0, fmt.Errorf("could not get total amount of rewards: %w", err)
	}
	return amount, nil
}
