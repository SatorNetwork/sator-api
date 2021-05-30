package challenge

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetChallengesByShowId endpoint.Endpoint
		GetChallengeById      endpoint.Endpoint
	}

	service interface {
		GetChallengeByID(ctx context.Context, id uuid.UUID) (interface{}, error)
		GetChallengesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error)
	}

	GetChallengeByIdRequest struct {
		ID string `json:"id" validate:"required,uuid"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetChallengeById: MakeGetChallengeByIdEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetChallengeById = mdw(e.GetChallengeById)
		}
	}

	return e
}

// MakeGetChallengeByIdEndpoint ...
func MakeGetChallengeByIdEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetChallengeByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
