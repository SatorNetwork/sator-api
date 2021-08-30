package referrals

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetMyReferralCode      endpoint.Endpoint
		StoreUserWithValidCode endpoint.Endpoint
	}

	service interface {
		GetMyReferralCode(ctx context.Context, uid uuid.UUID) (Data, error)
		StoreUserWithValidCode(ctx context.Context, uid uuid.UUID, code string) (bool, error)
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetMyReferralCode:      MakeGetMyReferralCodeEndpoint(s),
		StoreUserWithValidCode: MakeStoreUserWithValidCodeEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetMyReferralCode = mdw(e.GetMyReferralCode)
			e.StoreUserWithValidCode = mdw(e.StoreUserWithValidCode)
		}
	}

	return e
}

// MakeGetMyReferralCodeEndpoint ...
func MakeGetMyReferralCodeEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		resp, err := s.GetMyReferralCode(ctx, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeStoreUserWithValidCodeEndpoint ...
func MakeStoreUserWithValidCodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		resp, err := s.StoreUserWithValidCode(ctx, uid, request.(string))
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
