package wallet

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		wr walletRepository
		sc solana.Client
		// rw rewardsService

		satorAssetName  string
		solanaAssetName string

		walletDetailsURL        string // url template to get SOL & SAO wallet types details
		walletTransactionsURL   string // url template to get SOL & SAO wallet types transactions list
		rewardsWalletDetailsURL string // url template to get rewards wallet type details
		rewardsTransactionsURL  string // url template to get rewards wallet type transactions list
	}

	walletRepository interface {
		CreateWallet(ctx context.Context, arg repository.CreateWalletParams) (repository.Wallet, error)
		GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]repository.Wallet, error)
		GetWalletBySolanaAccountID(ctx context.Context, solanaAccountID uuid.UUID) (repository.Wallet, error)
		GetWalletByID(ctx context.Context, id uuid.UUID) (repository.Wallet, error)

		AddSolanaAccount(ctx context.Context, arg repository.AddSolanaAccountParams) (repository.SolanaAccount, error)
		GetSolanaAccountByID(ctx context.Context, id uuid.UUID) (repository.SolanaAccount, error)
		GetSolanaAccountByUserID(ctx context.Context, userID uuid.UUID) (repository.SolanaAccount, error)
	}

	// rewardsService interface {
	// 	GetTotalAmount(ctx context.Context, userID uuid.UUID) (float64, error)
	// }
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(wr walletRepository, sc solana.Client) *Service {
	if wr == nil {
		log.Fatalln("wallet repository is not set")
	}
	return &Service{
		wr: wr,
		sc: sc,
		// rw: rw,

		solanaAssetName: "SOL",
		satorAssetName:  "SAO",

		walletDetailsURL:        "wallets/%s",
		walletTransactionsURL:   "wallets/%s/transactions",
		rewardsWalletDetailsURL: "rewards/wallet/%s",
		rewardsTransactionsURL:  "rewards/wallet/%s/transactions",
	}
}

// GetWallets returns current user's wallets list with balance
func (s *Service) GetWallets(ctx context.Context, uid uuid.UUID) (Wallets, error) {
	wallets, err := s.wr.GetWalletsByUserID(ctx, uid)
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("could not get wallets list: %w", err)
	}

	result := make(Wallets, 0, len(wallets))
	for _, w := range wallets {
		wli := WalletsListItem{ID: w.ID.String(), Type: w.WalletType}

		switch w.WalletType {
		case TypePersonal:
			wli.GetDetailsURL = fmt.Sprintf(s.walletDetailsURL, w.ID.String())
			wli.GetTransactionsURL = fmt.Sprintf(s.walletTransactionsURL, w.ID.String())
		case TypeRewards:
			wli.GetDetailsURL = fmt.Sprintf(s.rewardsWalletDetailsURL, w.ID.String())
			wli.GetTransactionsURL = fmt.Sprintf(s.rewardsTransactionsURL, w.ID.String())
		}

		result = append(result, wli)
	}

	return result, nil
}

// GetWalletByID returns wallet details by wallet id
func (s *Service) GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (Wallet, error) {
	w, err := s.wr.GetWalletByID(ctx, walletID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return Wallet{}, fmt.Errorf("%w wallet", ErrNotFound)
		}
		return Wallet{}, fmt.Errorf("could not get wallet: %w", err)
	}

	if w.UserID != userID {
		return Wallet{}, fmt.Errorf("%w: you have no permissions to get this wallet", ErrForbidden)
	}

	sa, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
	if err != nil {
		if db.IsNotFoundError(err) {
			return Wallet{}, fmt.Errorf("%w solana account for this wallet", ErrNotFound)
		}
		return Wallet{}, fmt.Errorf("could not get solana account for this wallet: %w", err)
	}

	var balance []Balance

	balanceToken, err := s.sc.GetTokenAccountBalance(ctx, sa.PublicKey)
	if err != nil {
		return Wallet{}, fmt.Errorf("couldn't get token balance for this account: %w", err)
	}

	balanceSol, err := s.sc.GetAccountBalanceSOL(ctx, sa.PublicKey)
	if err != nil {
		return Wallet{}, fmt.Errorf("couldn't get solana balance for this account: %w", err)
	}

	balance = []Balance{
		{
			Currency: s.solanaAssetName,
			Amount:   balanceSol,
		},
		{
			Currency: "USD",
			Amount:   balanceSol*34.5 + balanceToken*1.25, // FIXME: setup currency rate
		},
	}

	return Wallet{
		ID:                   w.ID.String(),
		SolanaAccountAddress: sa.PublicKey,
		Actions: []Action{
			{
				Type: ActionSendTokens.String(),
				Name: ActionSendTokens.Name(),
				URL:  "",
			},
			{
				Type: ActionReceiveTokens.String(),
				Name: ActionReceiveTokens.Name(),
				URL:  "",
			},
		},
		Balance: balance,
	}, nil
}

// CreateWallet creates wallet for user with specified id.
func (s *Service) CreateWallet(ctx context.Context, userID uuid.UUID) error {
	acc := s.sc.NewAccount()

	txHash, err := s.sc.InitAccountToUseAsset(
		ctx, acc,
	)
	if err != nil {
		return fmt.Errorf("could not init token holder account: %w", err)
	}
	log.Printf("init token holder account transaction: %s", txHash)

	sacc, err := s.wr.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		PublicKey:  acc.PublicKey.ToBase58(),
		PrivateKey: acc.PrivateKey,
	})
	if err != nil {
		return fmt.Errorf("could not store solana account: %w", err)
	}

	if _, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:          userID,
		SolanaAccountID: sacc.ID,
		WalletType:      TypePersonal,
		Sort:            1,
	}); err != nil {
		return fmt.Errorf("could not new SAO wallet for user with id=%s: %w", userID.String(), err)
	}

	if _, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:     userID,
		WalletType: TypeRewards,
		Sort:       2,
	}); err != nil {
		return fmt.Errorf("could not new rewards wallet for user with id=%s: %w", userID.String(), err)
	}

	return nil
}

// WithdrawRewards convert rewards into sator tokens
func (s *Service) WithdrawRewards(ctx context.Context, userID uuid.UUID, amount float64) (tx string, err error) {
	user, err := s.wr.GetSolanaAccountByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("could not get user token account: %w", err)
	}

	// sends token
	for i := 0; i < 5; i++ {
		if tx, err = s.sc.SendAssets(
			ctx,
			s.sc.AccountFromPrivatekey(s.sc.Issuer.PrivateKey),
			user.PublicKey,
			amount,
		); err != nil {
			log.Println(err)
			time.Sleep(time.Second * 10)
		} else {
			log.Printf("user %s: successful transaction: rewards withdraw: %s", userID.String(), tx)
			break
		}
	}

	return tx, nil
}

// Bootstrap for usage in development mode only
func (s *Service) Bootstrap(ctx context.Context) error {
	if _, err := s.wr.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		PublicKey:  s.sc.FeePayer.PublicKey.ToBase58(),
		PrivateKey: s.sc.FeePayer.PrivateKey,
	}); err != nil {
		return fmt.Errorf("could not store issuer solana account: %w", err)
	}

	for i := 0; i < 5; i++ {
		tx, err := s.sc.RequestAirdrop(ctx, s.sc.FeePayer.PublicKey.ToBase58(), 10)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("airdrop for %s: %s", s.sc.FeePayer.PublicKey.ToBase58(), tx)
		break
	}

	time.Sleep(time.Second * 15)

	if _, err := s.wr.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		PublicKey:  s.sc.Asset.PublicKey.ToBase58(),
		PrivateKey: s.sc.Asset.PrivateKey,
	}); err != nil {
		return fmt.Errorf("could not store asset solana account: %w", err)
	}

	if tx, err := s.sc.CreateAsset(ctx); err != nil {
		return err
	} else {
		log.Printf("convert account (%s) to asset: %s", s.sc.Asset.PublicKey.ToBase58(), tx)
	}

	time.Sleep(time.Second * 15)

	if _, err := s.wr.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		PublicKey:  s.sc.Issuer.PublicKey.ToBase58(),
		PrivateKey: s.sc.Issuer.PrivateKey,
	}); err != nil {
		return fmt.Errorf("could not store issuer solana account: %w", err)
	}

	if tx, err := s.sc.InitAccountToUseAsset(
		ctx,
		s.sc.AccountFromPrivatekey(s.sc.Issuer.PrivateKey),
	); err != nil {
		return err
	} else {
		log.Printf("init issuer account (%s) to user asset: %s", s.sc.Issuer.PublicKey.ToBase58(), tx)
	}

	time.Sleep(time.Second * 15)

	if tx, err := s.sc.IssueAsset(
		ctx,
		s.sc.AccountFromPrivatekey(s.sc.Issuer.PrivateKey),
		1000000,
	); err != nil {
		return err
	} else {
		log.Printf("issue asset (%s): %s", s.sc.Issuer.PublicKey.ToBase58(), tx)
	}

	return nil
}

// GetListTransactionsByWalletID returns list of all transactions of specific wallet.
func (s *Service) GetListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID, limit, offset int32) (_ Transactions, err error) {
	transactions, err := s.getListTransactionsByWalletID(ctx, userID, walletID)
	if err != nil {
		return Transactions{}, err
	}

	transactionsPaginated := paginateTransactions(transactions, int(offset), int(limit))
	return transactionsPaginated, nil
}

// getListTransactionsByWalletID returns list of all transactions of specific wallet.
func (s *Service) getListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID) (Transactions, error) {
	wallet, err := s.wr.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	if wallet.UserID != userID {
		return nil, ErrForbidden
	}

	solanaAcc, err := s.wr.GetSolanaAccountByID(ctx, wallet.SolanaAccountID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.sc.GetTransactions(ctx, solanaAcc.PublicKey)
	if err != nil {
		return nil, err
	}

	txList := make(Transactions, 0, len(transactions))
	for _, tx := range transactions {
		txList = append(txList, castSolanaTxToTransaction(tx))
	}

	return txList, nil
}

// Transfer sends transaction from one account to another.
func (s *Service) Transfer(ctx context.Context, senderPrivateKey, recipientPK string, amount float64) (tx string, err error) {
	senderAcc := s.sc.AccountFromPrivatekey([]byte(senderPrivateKey))

	for i := 0; i < 5; i++ {
		if tx, err = s.sc.SendAssets(
			ctx,
			senderAcc,
			recipientPK,
			amount,
		); err != nil {
			log.Println(err)
			time.Sleep(time.Second * 10)
		} else {
			log.Printf("successful transaction: %s", tx)
			break
		}
	}

	return tx, nil
}

func castSolanaTxToTransaction(tx solana.ConfirmedTransactionResponse) Transaction {
	return Transaction{
		TxHash:    tx.TxHash,
		Amount:    tx.Amount,
		CreatedAt: tx.CreatedAt.Format(time.RFC3339),
	}
}

func paginateTransactions(transactions Transactions, offset, limit int) Transactions {
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt < (transactions[j].CreatedAt)
	})

	if offset > len(transactions) {
		offset = len(transactions)
	}

	end := offset + limit
	if end > len(transactions) {
		end = len(transactions)
	}

	return transactions[offset:end]
}
