package interfaces

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type RewardsService interface {
	AddDepositTransaction(ctx context.Context, userID, relationID uuid.UUID, relationType string, amount float64) error
}