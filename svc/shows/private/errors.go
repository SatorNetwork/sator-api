package private

import "errors"

// Predefined package errors
var (
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyReviewed  = errors.New("you have already written a review")
	ErrMaxClaps         = errors.New("the maximum possible number of applause for one show has been exceeded")
)
