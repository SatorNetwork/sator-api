package sumsub

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"

	"github.com/SatorNetwork/sator-api/lib/jwt"
)

// For a type to be a KYCStatus object, it must just have a GetKYCStatus method that returns user kyc status.
type kycStatus func(ctx context.Context, uid uuid.UUID) (string, error)

func KYCStatusMdw(keyFunc kycStatus, skipFunc func() bool) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if skipFunc != nil && skipFunc() {
				return next(ctx, request)
			}

			uid, err := jwt.UserIDFromContext(ctx)
			if err != nil {
				return nil, fmt.Errorf("could not get user id: %w", err)
			}

			s, err := keyFunc(ctx, uid)
			if err != nil {
				return nil, fmt.Errorf("could not check user kyc status: %w", err)
			}

			switch s {
			case KYCStatusApproved:
				return next(ctx, request)
			case KYCStatusRejected:
				return nil, ErrKYCUserIsDisabled
			case KYCStatusInProgress:
				return nil, ErrKYCInProgress
			case KYCStatusRetry:
				return nil, ErrKYCNeeded
			case KYCStatusInit:
				return nil, ErrKYCRequiredDocs
			default:
				return nil, ErrKYCNeeded
			}
		}
	}
}
