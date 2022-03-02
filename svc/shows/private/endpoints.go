package private

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/SatorNetwork/sator-api/internal/utils"
	"github.com/SatorNetwork/sator-api/internal/validator"
)

type (
	Endpoints struct {
		GetShows endpoint.Endpoint
	}

	service interface {
		GetShows(ctx context.Context, page, itemsPerPage int32) (interface{}, error)
	}

	GetShowsRequest struct {
		utils.PaginationRequest
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

func MakeGetShowsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetShowsRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetShows(ctx, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
