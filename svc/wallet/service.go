package wallet

import (
	"context"
	"log"

	"github.com/SatorNetwork/sator-api/svc/wallet/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		wr walletRepository
	}

	// WalletsBalance balance accross all wallets
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
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(wr walletRepository) *Service {
	if wr == nil {
		log.Fatalln("wallet repository is not set")
	}
	return &Service{wr: wr}
}

// GetBalance returns current user's balance
// TODO: take balance from solana
func (s *Service) GetBalance(ctx context.Context, uid uuid.UUID) (interface{}, error) {
	return WalletsBalance{
		"sao": Balance{
			Amount:   302,
			Currency: "SAO",
		},
		"usd": Balance{
			Amount:   2541.29,
			Currency: "USD",
		},
	}, nil
}
