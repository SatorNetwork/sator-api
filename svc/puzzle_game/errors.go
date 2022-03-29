package puzzle_game

import "errors"

// Predefined package errors
var (
	ErrNotFound         = errors.New("not found")
	ErrForbidden        = errors.New("forbidden")
	ErrInvalidParameter = errors.New("invalid parameter") // ErrInvalidParameter indicates that passed invalid parameter.
)
