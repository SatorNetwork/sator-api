package wallet

import "errors"

// Predefined package error
var (
	ErrInvalidParameter    = errors.New("invalid parameter") // ErrInvalidParameter indicates that passed invalid parameter.
	ErrNotFound            = errors.New("not found")
	ErrForbidden           = errors.New("forbidden")
	ErrTokenHolderBalance  = errors.New("the rewards pool for this period is out. Come back next time to claim rewards")
	ErrMinimalAmountToSend = errors.New("minimal amount to send")
	ErrNotEnoughBalance    = errors.New("not enough balance amount")
	ErrCouldNotLock        = errors.New("could not lock tokens, try again later")
	ErrTransactionFailed   = errors.New("transaction failed")
	ErrFraudDetection      = errors.New("fraud detection")
	ErrTooManyRequests     = errors.New("too many requests, please try again later")
)
