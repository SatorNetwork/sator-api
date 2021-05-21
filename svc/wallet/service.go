package wallet

import (
	"context"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
	}

	// WalletsBalance balance accross all wallets
	WalletsBalance map[string]Balance

	// Balance struct
	Balance struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService() *Service {
	return &Service{}
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
			Amount:   2541,
			Currency: "USD",
		},
	}, nil
}
