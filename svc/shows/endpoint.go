package shows

import (
	"context"

	"github.com/SatorNetwork/sator-api/internal/validator"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetShows endpoint.Endpoint
	}

	service interface {
		GetShows(ctx context.Context, page int32) (interface{}, error)
	}

	// GetShowsRequest struct
	GetShowsRequest struct {
		Page int32 `json:"page" validate:"required"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetShows: MakeGetShowsEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetShows = mdw(e.GetShows)
		}
	}

	return e
}

// MakeGetShowsEndpoint ...
func MakeGetShowsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetShowsRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetShows(ctx, req.Page)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
