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
	wallet, err := s.wr.GetWalletByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("could not get wallet by id: %w", err)
	}

	amount, err := s.sc.GetBalance(ctx, wallet.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("could not get balance: %w", err)
	}

	return WalletsBalance{
		"sao": Balance{
			Amount:   float64(amount),
			Currency: "SAO",
		},
		"usd": Balance{
			Amount:   float64(amount) * 25,
			Currency: "USD",
		},
	}, nil
}

// CreateWallet creates wallet for user with specified id
func (s *Service) CreateWallet(ctx context.Context, userID uuid.UUID, publicKey string, privateKey []byte) (repository.Wallet, error) {
	wallet, err := s.wr.CreateWallet(ctx, repository.CreateWalletParams{
		UserID:     userID,
		PublicKey:  publicKey,
		PrivateKey: string(privateKey),
	})
	if err != nil {
		return repository.Wallet{}, err
	}

	return wallet, nil
}
