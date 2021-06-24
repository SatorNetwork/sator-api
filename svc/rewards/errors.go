package rewards

import "errors"

var (
	// ErrRewardsAlreadyClaimed indicated that all rewards already claimed.
	ErrRewardsAlreadyClaimed = errors.New("you have already claimed all rewards")

	// ErrInvalidParameter indicates that passed invalid parameter.
	ErrInvalidParameter = errors.New("invalid parameter")
)
