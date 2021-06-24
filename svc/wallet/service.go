package wallet

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"

	"github.com/google/uuid"
	"github.com/portto/solana-go-sdk/types"
)

type (
	// Service struct
	Service struct {
		wr walletRepository
		sc solanaClient
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
		GetSolanaAccountByType(ctx context.Context, accountType string) (repository.SolanaAccount, error)
		GetSolanaAccountTypeByPublicKey(ctx context.Context, publicKey string) (string, error)
		GetSolanaAccountByUserIDAndType(ctx context.Context, arg repository.GetSolanaAccountByUserIDAndTypeParams) (repository.SolanaAccount, error)
	}

	solanaClient interface {
		GetAccountBalanceSOL(ctx context.Context, accPubKey string) (float64, error)
		GetTokenAccountBalance(ctx context.Context, accPubKey string) (float64, error)
		NewAccount() types.Account
		RequestAirdrop(ctx context.Context, pubKey string, amount float64) (string, error)
		AccountFromPrivatekey(pk []byte) types.Account
		InitAccountToUseAsset(ctx context.Context, feePayer, issuer, asset, initAcc types.Account) (string, error)
		SendAssets(ctx context.Context, feePayer, issuer, asset, sender types.Account, recipientAddr string, amount float64) (string, error)
		CreateAsset(ctx context.Context, feePayer, issuer, asset types.Account) (string, error)
		IssueAsset(ctx context.Context, feePayer, issuer, asset, dest types.Account, amount float64) (string, error)
		GetTransactions(ctx context.Context, publicKey string) ([]solana.ConfirmedTransactionResponse, error)
	}

	// rewardsService interface {
	// 	GetTotalAmount(ctx context.Context, userID uuid.UUID) (float64, error)
	// }
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(wr walletRepository, sc solanaClient) *Service {
	if wr == nil {
		log.Fatalln("wallet repository is not set")
	}
	if sc == nil {
		log.Fatalln("solana client is not set")
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
		case WalletTypeSolana, WalletTypeSator:
			wli.GetDetailsURL = fmt.Sprintf(s.walletDetailsURL, w.ID.String())
			wli.GetTransactionsURL = fmt.Sprintf(s.walletTransactionsURL, w.ID.String())
		case WalletTypeRewards:
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

	switch sa.AccountType {
	case TokenAccount.String():
		if bal, err := s.sc.GetTokenAccountBalance(ctx, sa.PublicKey); err == nil {
			balance = []Balance{
				{
					Currency: s.satorAssetName,
					Amount:   bal,
				},
				{
					Currency: "USD",
					Amount:   bal * 1.25, // FIXME: setup currency rate
				},
			}
		}
	case GeneralAccount.String():
		if bal, err := s.sc.GetAccountBalanceSOL(ctx, sa.PublicKey); err == nil {
			balance = []Balance{
				{
					Currency: s.solanaAssetName,
					Amount:   bal,
				},
				{
					Currency: "USD",
					Amount:   bal * 34.5, // FIXME: setup currency rate
				},
			}
		}
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
	feePayer, err := s.wr.GetSolanaAccountByType(ctx, FeePayerAccount.String())
	if err != nil {
		return fmt.Errorf("could not get fee payer account: %w", err)
	}
	issuer, err := s.wr.GetSolanaAccountByType(ctx, IssuerAccount.String())
	if err != nil {
		return fmt.Errorf("could not get issuer account: %w", err)
	}
	asset, err := s.wr.GetSolanaAccountByType(ctx, AssetAccount.String())
	if err != nil {
		return fmt.Errorf("could not get asset account: %w", err)
	}

	acc := s.sc.NewAccount()

	txHash, err := s.sc.InitAccountToUseAsset(
		ctx,
		s.sc.AccountFromPrivatekey(feePayer.PrivateKey),
		s.sc.AccountFromPrivatekey(issuer.PrivateKey),
		s.sc.AccountFromPrivatekey(asset.PrivateKey),
		acc,
	)
	if err != nil {
		return fmt.Errorf("could not init token holder account: %w", err)
	}
	log.Printf("init token holder account transaction: %s", txHash)

	sacc, err := s.wr.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		AccountType: TokenAccount.String(),
		PublicKey:   acc.PublicKey.ToBase58(),
		PrivateKey:  acc.PrivateKey,
	})
	if err != nil {
		return fmt.Errorf("could not store solana account: %w", err)
	}

	if _, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:          userID,
		SolanaAccountID: sacc.ID,
		WalletType:      WalletTypeSator,
		Sort:            1,
	}); err != nil {
		return fmt.Errorf("could not new SAO wallet for user with id=%s: %w", userID.String(), err)
	}

	if _, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:     userID,
		WalletType: WalletTypeRewards,
		Sort:       2,
	}); err != nil {
		return fmt.Errorf("could not new rewards wallet for user with id=%s: %w", userID.String(), err)
	}

	return nil
}

// WithdrawRewards convert rewards into sator tokens
func (s *Service) WithdrawRewards(ctx context.Context, userID uuid.UUID, amount float64) (tx string, err error) {
	feePayer, err := s.wr.GetSolanaAccountByType(ctx, FeePayerAccount.String())
	if err != nil {
		return "", fmt.Errorf("could not get fee payer account: %w", err)
	}
	issuer, err := s.wr.GetSolanaAccountByType(ctx, IssuerAccount.String())
	if err != nil {
		return "", fmt.Errorf("could not get issuer account: %w", err)
	}
	asset, err := s.wr.GetSolanaAccountByType(ctx, AssetAccount.String())
	if err != nil {
		return "", fmt.Errorf("could not get asset account: %w", err)
	}
	user, err := s.wr.GetSolanaAccountByUserIDAndType(ctx, repository.GetSolanaAccountByUserIDAndTypeParams{
		UserID:     userID,
		WalletType: WalletTypeSator,
	})
	if err != nil {
		return "", fmt.Errorf("could not get user token account: %w", err)
	}

	// sends token
	for i := 0; i < 5; i++ {
		if tx, err = s.sc.SendAssets(
			ctx,
			s.sc.AccountFromPrivatekey(feePayer.PrivateKey),
			s.sc.AccountFromPrivatekey(issuer.PrivateKey),
			s.sc.AccountFromPrivatekey(asset.PrivateKey),
			s.sc.AccountFromPrivatekey(issuer.PrivateKey),
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
	feePayer := s.sc.NewAccount()
	issuer := s.sc.NewAccount()
	asset := s.sc.NewAccount()

	if _, err := s.wr.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		AccountType: FeePayerAccount.String(),
		PublicKey:   feePayer.PublicKey.ToBase58(),
		PrivateKey:  feePayer.PrivateKey,
	}); err != nil {
		return fmt.Errorf("could not store issuer solana account: %w", err)
	}

	for i := 0; i < 5; i++ {
		tx, err := s.sc.RequestAirdrop(ctx, feePayer.PublicKey.ToBase58(), 10)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("airdrop for %s: %s", feePayer.PublicKey.ToBase58(), tx)
		break
	}

	time.Sleep(time.Second * 15)

	if _, err := s.wr.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		AccountType: AssetAccount.String(),
		PublicKey:   asset.PublicKey.ToBase58(),
		PrivateKey:  asset.PrivateKey,
	}); err != nil {
		return fmt.Errorf("could not store asset solana account: %w", err)
	}

	if tx, err := s.sc.CreateAsset(
		ctx,
		s.sc.AccountFromPrivatekey(feePayer.PrivateKey),
		s.sc.AccountFromPrivatekey(issuer.PrivateKey),
		s.sc.AccountFromPrivatekey(asset.PrivateKey),
	); err != nil {
		return err
	} else {
		log.Printf("convert account (%s) to asset: %s", asset.PublicKey.ToBase58(), tx)
	}

	time.Sleep(time.Second * 15)

	if _, err := s.wr.AddSolanaAccount(ctx, repository.AddSolanaAccountParams{
		AccountType: IssuerAccount.String(),
		PublicKey:   issuer.PublicKey.ToBase58(),
		PrivateKey:  issuer.PrivateKey,
	}); err != nil {
		return fmt.Errorf("could not store issuer solana account: %w", err)
	}

	if tx, err := s.sc.InitAccountToUseAsset(
		ctx,
		s.sc.AccountFromPrivatekey(feePayer.PrivateKey),
		s.sc.AccountFromPrivatekey(issuer.PrivateKey),
		s.sc.AccountFromPrivatekey(asset.PrivateKey),
		s.sc.AccountFromPrivatekey(issuer.PrivateKey),
	); err != nil {
		return err
	} else {
		log.Printf("init issuer account (%s) to user asset: %s", issuer.PublicKey.ToBase58(), tx)
	}

	time.Sleep(time.Second * 15)

	if tx, err := s.sc.IssueAsset(
		ctx,
		s.sc.AccountFromPrivatekey(feePayer.PrivateKey),
		s.sc.AccountFromPrivatekey(issuer.PrivateKey),
		s.sc.AccountFromPrivatekey(asset.PrivateKey),
		s.sc.AccountFromPrivatekey(issuer.PrivateKey),
		1000000,
	); err != nil {
		return err
	} else {
		log.Printf("issue asset (%s): %s", issuer.PublicKey.ToBase58(), tx)
	}

	return nil
}

// GetListTransactionsByWalletID returns list of all transactions of specific wallet.
func (s *Service) GetListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID) (_ Transactions, err error) {
	return s.getListTransactionsByWalletID(ctx, userID, walletID)
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

	senderAccType, err := s.wr.GetSolanaAccountByType(ctx, senderAcc.PublicKey.ToBase58())
	if err != nil {
		return "", fmt.Errorf("could not get fee payer account: %w", err)
	}

	recipientAccType, err := s.wr.GetSolanaAccountByType(ctx, recipientPK)
	if err != nil {
		return "", fmt.Errorf("could not get fee payer account: %w", err)
	}

	if senderAccType.AccountType != recipientAccType.AccountType {
		return "", fmt.Errorf("accounts have different types, transaction impossible")
	}

	feePayer, err := s.wr.GetSolanaAccountByType(ctx, FeePayerAccount.String())
	if err != nil {
		return "", fmt.Errorf("could not get fee payer account: %w", err)
	}
	issuer, err := s.wr.GetSolanaAccountByType(ctx, IssuerAccount.String())
	if err != nil {
		return "", fmt.Errorf("could not get issuer account: %w", err)
	}
	asset, err := s.wr.GetSolanaAccountByType(ctx, AssetAccount.String())
	if err != nil {
		return "", fmt.Errorf("could not get asset account: %w", err)
	}

	for i := 0; i < 5; i++ {
		if tx, err = s.sc.SendAssets(
			ctx,
			s.sc.AccountFromPrivatekey(feePayer.PrivateKey),
			s.sc.AccountFromPrivatekey(issuer.PrivateKey),
			s.sc.AccountFromPrivatekey(asset.PrivateKey),
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
