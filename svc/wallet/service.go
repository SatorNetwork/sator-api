package wallet

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/svc/wallet/repository"

	"github.com/google/uuid"
	"github.com/portto/solana-go-sdk/types"
)

type (
	// Service struct
	Service struct {
		wr              walletRepository
		sc              solanaClient
		rw              rewardsService
		satorAssetName  string
		solanaAssetName string
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

	rewardsService interface {
		GetTotalAmount(ctx context.Context, userID uuid.UUID) (float64, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(wr walletRepository, sc solanaClient, rw rewardsService) *Service {
	if wr == nil {
		log.Fatalln("wallet repository is not set")
	}
	if sc == nil {
		log.Fatalln("solana client is not set")
	}
	return &Service{
		wr:              wr,
		sc:              sc,
		rw:              rw,
		solanaAssetName: "SOL",
		satorAssetName:  "SAO",
	}
}

// GetWallets returns current user's wallets list with balance
func (s *Service) GetWallets(ctx context.Context, uid uuid.UUID) (interface{}, error) {
	balance, err := s.getWalletsBalance(ctx, uid)
	if err != nil {
		return nil, err
	}

	rewAmount, err := s.rw.GetTotalAmount(ctx, uid)
	if err != nil {
		rewAmount = 0
	}

	balance = append(balance, Balance{
		Currency: "rewards",
		Amount:   rewAmount,
	})

	return balance, nil
}

// GetBalanceWithRewards returns current user's balance
func (s *Service) GetBalanceWithRewards(ctx context.Context, uid uuid.UUID) (interface{}, error) {
	balance, err := s.getWalletsBalance(ctx, uid)
	if err != nil {
		return nil, err
	}

	rewAmount, err := s.rw.GetTotalAmount(ctx, uid)
	if err != nil {
		rewAmount = 0
	}

	balance = append(balance, Balance{
		Currency: "rewards",
		Amount:   rewAmount,
	})

	return balance, nil
}

// getWalletsBalance returns wallet's solana account address, currency and amount.
func (s *Service) getWalletsBalance(ctx context.Context, userID uuid.UUID) ([]Balance, error) {
	wallets, err := s.wr.GetWalletsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("could not get wallets of current user [%s]: %w", userID.String(), err)
	}

	result := make([]Balance, 0, len(wallets))
	for _, w := range wallets {
		sa, err := s.wr.GetSolanaAccountByID(ctx, w.SolanaAccountID)
		if err != nil {
			log.Printf("could not get solana account with id=%s: %v", w.SolanaAccountID.String(), err)
			continue
		}

		var currency string
		var amount float64

		switch sa.AccountType {
		case TokenAccount.String():
			currency = s.satorAssetName
			if bal, err := s.sc.GetTokenAccountBalance(ctx, sa.PublicKey); err == nil {
				amount = bal
			}
		case GeneralAccount.String():
			currency = s.solanaAssetName
			if bal, err := s.sc.GetAccountBalanceSOL(ctx, sa.PublicKey); err == nil {
				amount = bal
			}
		}

		result = append(result, Balance{
			// SolanaAccountAddress: sa.PublicKey,
			Currency: currency,
			Amount:   amount,
		})
	}

	return result, nil
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
	}); err != nil {
		return fmt.Errorf("could not new SAO wallet for user with id=%s: %w", userID.String(), err)
	}

	if _, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:     userID,
		WalletType: WalletTypeRewards,
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
		UserID:      userID,
		AccountType: TokenAccount.String(),
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
func (s *Service) GetListTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) (_ interface{}, err error) {
	return s.getListTransactionsByWalletID(ctx, walletID)
}

// getListTransactionsByWalletID returns list of all transactions of specific wallet.
func (s *Service) getListTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]Transaction, error) {
	wallet, err := s.wr.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	solanaAcc, err := s.wr.GetSolanaAccountByID(ctx, wallet.SolanaAccountID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.sc.GetTransactions(ctx, solanaAcc.PublicKey)
	if err != nil {
		return nil, err
	}

	txList := make([]Transaction, 0, len(transactions))
	for _, tx := range transactions {
		txList = append(txList, castSolanaTxToTransaction(tx))
	}

	return txList, nil
}

// GetBalanceByUserID returns all user's wallets balance info's.
func (s *Service) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) ([]Balance, error) {
	return s.getWalletsBalance(ctx, userID)
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
		CreatedAt: tx.CreatedAt.String(),
	}
}
