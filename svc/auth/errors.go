package auth

import "errors"

// Predefined errors of the auth package
var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrOTPCode              = errors.New("invalid code")
	ErrNotFound             = errors.New("not found")
	ErrEmailAlreadyTaken    = errors.New("given email is already taken")
	ErrEmailAlreadyVerified = errors.New("your email address is already verified")

	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)
