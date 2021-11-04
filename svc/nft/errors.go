package nft

import "errors"

// Predefined package errors
var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrAlreadySold      = errors.New("reselling is not available in the current app version")
)
