package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type StakeLevels interface {
	GetMultiplier(ctx context.Context, userID uuid.UUID) (_ int32, err error)
}

type StaticStakeLevel struct{}

func (mock *StaticStakeLevel) GetMultiplier(ctx context.Context, userID uuid.UUID) (_ int32, err error) {
	return 1, nil
}
