package profile

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetProfile endpoint.Endpoint
	}

	service interface {
		GetProfileByUserID(ctx context.Context, userID uuid.UUID) (interface{}, error)
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GetProfile: MakeGetProfileEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetProfile = mdw(e.GetProfile)
		}
	}

	return e
}

// MakeGetProfileEndpoint ...
func MakeGetProfileEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		uid, ok := req.(uuid.UUID)
		if !ok {
			return nil, fmt.Errorf("could not get user profile id")
		}

		resp, err := s.GetProfileByUserID(ctx, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
