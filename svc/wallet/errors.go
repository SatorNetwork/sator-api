package wallet

import "errors"

// Predefined package error
var (
	ErrInvalidParameter       = errors.New("invalid parameter") // ErrInvalidParameter indicates that passed invalid parameter.
	ErrNotFound               = errors.New("not found")
	ErrForbidden              = errors.New("forbidden")
	ErrTokenHolderBalance     = errors.New("could not claim rewards, please try again later")
	ErrTemporarilyUnavailable = errors.New("wallet is temporarily unavailable")
)
