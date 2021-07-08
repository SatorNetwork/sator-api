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
		AddChallenge          endpoint.Endpoint
		DeleteChallengeByID   endpoint.Endpoint
		UpdateChallenge       endpoint.Endpoint
	}

	service interface {
		GetChallengeByID(ctx context.Context, id uuid.UUID) (interface{}, error)
		GetChallengesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error)
		AddChallenge(ctx context.Context, ch Challenge) error
		DeleteChallengeByID(ctx context.Context, id uuid.UUID) error
		UpdateChallenge(ctx context.Context, ch Challenge) error
	}

	// AddChallengeRequest struct
	AddChallengeRequest struct {
		ShowID          string `json:"show_id" validate:"required,uuid"`
		Title           string `json:"title" validate:"required,gt=0"`
		Description     string `json:"description"`
		PrizePool       string `json:"prize_pool" validate:"required,gt=0"`
		PlayersToStart  int32  `json:"players_to_start" validate:"required"`
		TimePerQuestion string `json:"time_per_question"`
	}

	// UpdateChallengeRequest struct
	UpdateChallengeRequest struct {
		ID              string `json:"id" validate:"required,uuid"`
		ShowID          string `json:"show_id" validate:"required,uuid"`
		Title           string `json:"title" validate:"required,gt=0"`
		Description     string `json:"description"`
		PrizePool       string `json:"prize_pool" validate:"required,gt=0"`
		PlayersToStart  int32  `json:"players_to_start" validate:"required"`
		TimePerQuestion string `json:"time_per_question"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetChallengeById:    MakeGetChallengeByIdEndpoint(s),
		AddChallenge:        MakeAddChallengeEndpoint(s, validateFunc),
		DeleteChallengeByID: MakeDeleteChallengeByIDEndpoint(s),
		UpdateChallenge:     MakeUpdateChallengeEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetChallengeById = mdw(e.GetChallengeById)
			e.AddChallenge = mdw(e.AddChallenge)
			e.DeleteChallengeByID = mdw(e.DeleteChallengeByID)
			e.UpdateChallenge = mdw(e.UpdateChallenge)
		}
	}

	return e
}

// MakeGetChallengeByIdEndpoint ...
func MakeGetChallengeByIdEndpoint(s service) endpoint.Endpoint {
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

// MakeAddChallengeEndpoint ...
func MakeAddChallengeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddChallengeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		err = s.AddChallenge(ctx, Challenge{
			ShowID:          showID,
			Title:           req.Title,
			Description:     req.Description,
			PrizePool:       req.PrizePool,
			Players:         req.PlayersToStart,
			TimePerQuestion: req.TimePerQuestion,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeDeleteChallengeByIDEndpoint ...
func MakeDeleteChallengeByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		err = s.DeleteChallengeByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%w challenge id: %v", ErrInvalidParameter, err)
		}

		return true, nil
	}
}

// MakeUpdateChallengeEndpoint ...
func MakeUpdateChallengeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateChallengeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		err = s.UpdateChallenge(ctx, Challenge{
			ID:              id,
			ShowID:          showID,
			Title:           req.Title,
			Description:     req.Description,
			PrizePool:       req.PrizePool,
			Players:         req.PlayersToStart,
			TimePerQuestion: req.TimePerQuestion,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}
