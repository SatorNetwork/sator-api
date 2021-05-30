package wallet

import (
	"context"
	"fmt"
	"log"

	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		wr walletRepository
		sc solanaClient
	}

	// WalletsBalance balance across all wallets
	WalletsBalance map[string]Balance

	// Balance struct
	Balance struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}

	walletRepository interface {
		CreateWallet(ctx context.Context, arg repository.CreateWalletParams) (repository.Wallet, error)
		GetWalletByAssetName(ctx context.Context, arg repository.GetWalletByAssetNameParams) (repository.Wallet, error)
		GetWalletByID(ctx context.Context, id uuid.UUID) (repository.Wallet, error)
		GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]repository.Wallet, error)
	}

	solanaClient interface {
		CreateAccount(ctx context.Context) (string, []byte, error)
		GetBalance(ctx context.Context, base58key string) (uint64, error)
	}
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
	return &Service{wr: wr, sc: sc}
}

// GetBalance returns current user's balance
func (s *Service) GetBalance(ctx context.Context, uid uuid.UUID) (interface{}, error) {
	wallets, err := s.wr.GetWalletsByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("could not get wallets of current user [%s]: %w", uid.String(), err)
	}

	var amount float64
	if len(wallets) > 0 {
		amount = s.getBalanceForWallet(ctx, wallets[0].PublicKey)
	}

	return WalletsBalance{
		"sao": Balance{
			Amount:   amount,
			Currency: "SAO",
		},
		"usd": Balance{
			Amount:   amount * 25,
			Currency: "USD",
		},
	}, nil
}

// CreateWallet creates wallet for user with specified id
func (s *Service) CreateWallet(ctx context.Context, userID uuid.UUID) (repository.Wallet, error) {
	publicKey, privateKey, err := s.sc.CreateAccount(ctx)
	if err != nil {
		return repository.Wallet{}, fmt.Errorf("could not create solana account: %w", err)
	}

	wallet, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:     userID,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})
	if err != nil {
		return repository.Wallet{}, err
	}

	return wallet, nil
}

func (s *Service) getBalanceForWallet(ctx context.Context, pubKey string) float64 {
	amount, err := s.sc.GetBalance(ctx, pubKey)
	if err != nil {
		log.Printf("could not get balance wir wallet %s: %v", pubKey, err)
		return 0
	}
	return toSol(amount)
}

func toSol(income uint64) float64 {
	return float64(income / 1000000000)
}
