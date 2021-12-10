package sumsub

import "errors"

// KYC possible errors
var (
	ErrNotFound          = errors.New("not found")
	ErrKYCNeeded         = errors.New("verification needed")
	ErrKYCInProgress     = errors.New("verification still in progress")
	ErrKYCUserIsDisabled = errors.New("your profile was disabled. Please contact support for details")
)
