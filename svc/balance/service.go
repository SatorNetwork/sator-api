package balance

import (
	"context"
	"log"

	"github.com/SatorNetwork/sator-api/svc/wallet"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		wallets walletsServiceClient
		rewards rewardsServiceClient
	}

	walletsServiceClient interface {
		GetWalletsListByUserID(ctx context.Context, userID uuid.UUID) (wallet.Wallets, error)
		GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (wallet.Wallet, error)
	}

	rewardsServiceClient interface {
		GetUserRewards(ctx context.Context, uid uuid.UUID) (float64, error)
	}

	// Balance struct
	Balance struct {
		Currency string  `json:"currency"`
		Amount   float64 `json:"amount"`
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(w walletsServiceClient, r rewardsServiceClient) *Service {
	return &Service{
		wallets: w,
		rewards: r,
	}
}

// GetAccountBalance returns user balance, including rewards
func (s *Service) GetAccountBalance(ctx context.Context, uid uuid.UUID) (interface{}, error) {
	walletsList, err := s.wallets.GetWalletsListByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}

	balance := make([]Balance, 0, len(walletsList))
	for _, w := range walletsList {
		switch w.Type {
		case wallet.WalletTypeSator, wallet.WalletTypeSolana:
			wlt, err := s.wallets.GetWalletByID(ctx, uid, uuid.MustParse(w.ID))
			if err != nil {
				log.Printf("get wallet with id=%s: %v", w.ID, err)
				continue
			}
			if len(wlt.Balance) > 0 {
				balance = append(balance, Balance{
					Currency: wlt.Balance[0].Currency,
					Amount:   wlt.Balance[0].Amount,
				})
			}
		case wallet.WalletTypeRewards:
			amount, err := s.rewards.GetUserRewards(ctx, uid)
			if err != nil {
				log.Printf("get user rewards: %v", err)
				continue
			}
			balance = append(balance, Balance{
				Currency: "unclaimed",
				Amount:   amount,
			})
		}
	}

	return balance, nil
}
