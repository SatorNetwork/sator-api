package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type WalletService interface {
	WithdrawRewards(ctx context.Context, userID uuid.UUID, amount float64) (string, error)
}
