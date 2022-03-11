//go:build !mock_solana

package client

import "errors"

// Predefined package errors
var (
	ErrATANotCreated = errors.New("associated token account does not exist")
)
