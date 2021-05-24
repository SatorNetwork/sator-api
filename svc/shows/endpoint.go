package shows

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
		GetShows          endpoint.Endpoint
		GetShowChallenges endpoint.Endpoint
	}

	service interface {
		GetShows(ctx context.Context, page, itemsPerPage int32) (interface{}, error)
		GetShowChallenges(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error)
	}

	// PaginationRequest struct
	PaginationRequest struct {
		Page         int32 `json:"page,omitempty" validate:"number,gte=0"`
		ItemsPerPage int32 `json:"items_per_page,omitempty" validate:"number,gte=0"`
	}

	// GetChallengesRequest struct
	GetChallengesRequest struct {
		ShowID string `json:"show_id" validate:"required,uuid"`
		PaginationRequest
	}
)

// Limit of items
func (r PaginationRequest) Limit() int32 {
	if r.ItemsPerPage > 0 {
		return r.ItemsPerPage
	}
	return 20
}

// Offset items
func (r PaginationRequest) Offset() int32 {
	if r.Page > 1 {
		return r.Page * r.Limit()
	}
	return 0
}

// MakeEndpoints ...
func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetShows:          MakeGetShowsEndpoint(s, validateFunc),
		GetShowChallenges: MakeGetShowChallengesEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetShows = mdw(e.GetShows)
			e.GetShowChallenges = mdw(e.GetShowChallenges)
		}
	}

	return e
}

// MakeGetShowsEndpoint ...
func MakeGetShowsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PaginationRequest)
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

// MakeGetShowChallengesEndpoint ...
func MakeGetShowChallengesEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetChallengesRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetShowChallenges(ctx, showID, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
