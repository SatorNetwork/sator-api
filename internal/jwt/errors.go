package jwt

import "errors"

// Predefined errors of jwt wrapper package
var (
	ErrUserIDEmpty       = errors.New("user id is empty")
	ErrJWTIDEmpty        = errors.New("jwt id is empty")
	ErrJWTSubjectEmpty   = errors.New("jwt subject is empty")
	ErrInvalidJWTClaims  = errors.New("invalid jwt claims")
	ErrInvalidJWTSubject = errors.New("invalid jwt subject")
)
