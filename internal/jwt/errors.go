package jwt

import "errors"

// Predefined errors of jwt wrapper package
var (
	ErrUserIDEmpty      = errors.New("user id is empty")
	ErrJWTIDEmpty       = errors.New("jwt id is empty")
	ErrInvalidJWTClaims = errors.New("invalid jwt claims")
)
