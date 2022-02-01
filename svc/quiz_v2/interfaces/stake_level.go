package interfaces

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type StakeLevels interface {
	GetMultiplier(ctx context.Context, userID uuid.UUID) (_ int32, err error)
}

type StaticStakeLevel struct{}

func (mock *StaticStakeLevel) GetMultiplier(ctx context.Context, userID uuid.UUID) (_ int32, err error) {
	return 1, nil
}