package nft

import "errors"

// Predefined package errors
var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrAlreadySold      = errors.New("this nft sold out")
)
