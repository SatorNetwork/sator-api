package gapi

import "errors"

// Predefined package errors
var (
	ErrCouldNotSignResponse    = errors.New("could not sign response")
	ErrCouldNotVerifySignature = errors.New("could not verify signature")
)
