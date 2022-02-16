package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type RewardsService interface {
	AddDepositTransaction(ctx context.Context, userID, relationID uuid.UUID, relationType string, amount float64) error
}
