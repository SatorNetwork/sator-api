package auth

import "errors"

// Predefined errors of the auth package
var (
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrOTPCode                  = errors.New("invalid code")
	ErrNotFound                 = errors.New("not found")
	ErrEmailAlreadyTaken        = errors.New("given email is already taken")
	ErrEmailAlreadyVerified     = errors.New("your email address is already verified")
	ErrMissedUserID             = errors.New("missed user id")
	ErrUserIsDisabled           = errors.New("your profile was disabled. Please contact support for details")
	ErrRestrictedEmailDomain    = errors.New("please use real email address, or contact administrator")
	ErrInvalidEmailFormat       = errors.New("email must be a valid email address")
	ErrEmptyDeviceID            = errors.New("the current version of the application is outdated please update to the latest version")
	ErrInvalidParameter         = errors.New("invalid parameter")
	ErrPublicKeyIsNotRegistered = errors.New("public key is not registered")

	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)
