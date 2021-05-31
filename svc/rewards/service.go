package rewards

import (
	"context"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

type (
	// Service struct
	Service struct {
		ws walletService
	}

	Winner struct {
		UserID uuid.UUID
		Points int
	}

	walletService interface {
		SendToWallet(ctx context.Context, userID uuid.UUID, amount float64) (string, error)
	}
)

func (service *Service) DistributeRewards(ctx context.Context, prizePool float64, winners []Winner) (err error) {
	var totalPoints int

	for _, winner := range winners {
		totalPoints += winner.Points
	}
	pointCost := prizePool / float64(totalPoints)

	for _, winner := range winners {
		_, sendErr := service.ws.SendToWallet(ctx, winner.UserID, pointCost*float64(winner.Points))
		if sendErr != nil {
			err = errs.Combine(err, sendErr)
			continue
		}
	}

	return err
}
