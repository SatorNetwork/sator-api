package balance

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetAccountBalance endpoint.Endpoint
	}

	service interface {
		GetAccountBalance(ctx context.Context, uid uuid.UUID) (interface{}, error)
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GetAccountBalance: MakeGetAccountBalanceEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetAccountBalance = mdw(e.GetAccountBalance)
		}
	}

	return e
}

// MakeGetAccountBalanceEndpoint ...
func MakeGetAccountBalanceEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		resp, err := s.GetAccountBalance(ctx, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
