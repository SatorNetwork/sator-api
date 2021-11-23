package wallet

import "errors"

// Predefined package error
var (
	ErrInvalidParameter   = errors.New("invalid parameter") // ErrInvalidParameter indicates that passed invalid parameter.
	ErrNotFound           = errors.New("not found")
	ErrForbidden          = errors.New("forbidden")
	ErrTokenHolderBalance = errors.New("the rewards pool for this period is out. Come back next time to claim rewards")
	ErrNotEnoughBalance   = errors.New("minimal amount to send")
)
