package repository

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/sumsub"
	"github.com/SatorNetwork/sator-api/svc/auth"
)

func KYCStatusMdw(keyFunc auth.KYCStatus) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			uid, err := jwt.UserIDFromContext(ctx)
			if err != nil {
				return nil, fmt.Errorf("could not get user id: %w", err)
			}

			s, err := keyFunc(uid)
			if err != nil {
				return nil, fmt.Errorf("could not get user id: %w", err)
			}

			switch s {
			case sumsub.KYCStatusApproved:
				return next(ctx, request)
			case sumsub.KYCStatusRejected:
				return nil, auth.ErrUserIsDisabled
			case sumsub.KYCStatusInProgress:
				return nil, auth.ErrKYCInProgress
			case sumsub.KYCStatusRetry:
				return nil, auth.ErrKYCNeeded
			default:
				return nil, auth.ErrKYCNeeded
			}
		}
	}
}
