package gapi

import "errors"

// Predefined package errors
var (
	ErrCouldNotStartGame    = errors.New("could not start game")
	ErrCouldNotFinishGame   = errors.New("could not finish game")
	ErrCouldNotSignResponse = errors.New("could not sign response")
)
