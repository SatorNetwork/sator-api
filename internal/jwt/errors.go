package jwt

import "errors"

// Predefined errors of jwt wrapper package
var (
	ErrUserIDEmpty       = errors.New("user id is empty")
	ErrJWTIDEmpty        = errors.New("jwt id is empty")
	ErrJWTSubjectEmpty   = errors.New("jwt subject is empty")
	ErrInvalidJWTClaims  = errors.New("invalid jwt claims")
	ErrInvalidJWTSubject = errors.New("invalid jwt subject")

	ErrUserIsDisabled = errors.New("your profile was disabled. Please contact support for details")
	ErrMissedUserID   = errors.New("missed user id")
)
