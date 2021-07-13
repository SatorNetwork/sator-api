package challenge

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
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

		GetVerificationQuestionByEpisodeID endpoint.Endpoint
		CheckVerificationQuestionAnswer    endpoint.Endpoint
		VerifyUserAccessToEpisode          endpoint.Endpoint
	}

	service interface {
		GetChallengeByID(ctx context.Context, id uuid.UUID) (interface{}, error)
		GetChallengesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error)
		AddChallenge(ctx context.Context, ch Challenge) (Challenge, error)
		DeleteChallengeByID(ctx context.Context, id uuid.UUID) error
		UpdateChallenge(ctx context.Context, ch Challenge) error

		// FIXME: needs refactoring!
		GetVerificationQuestionByEpisodeID(ctx context.Context, episodeID uuid.UUID) (interface{}, error)
		CheckVerificationQuestionAnswer(ctx context.Context, qid, aid uuid.UUID) (interface{}, error)
		VerifyUserAccessToEpisode(ctx context.Context, uid, eid uuid.UUID) (interface{}, error)
	}

	// AddChallengeRequest struct
	AddChallengeRequest struct {
		ShowID             string  `json:"show_id" validate:"required,uuid"`
		Title              string  `json:"title" validate:"required,gt=0"`
		Description        string  `json:"description"`
		PrizePoolAmount    float64 `json:"prize_pool_amount" validate:"required,gt=0"`
		PlayersToStart     int32   `json:"players_to_start" validate:"required"`
		TimePerQuestionSec int64   `json:"time_per_question_sec"`
		EpisodeID          string  `json:"episode_id" validate:"uuid"`
		Kind               int32   `json:"kind"`
	}

	// UpdateChallengeRequest struct
	UpdateChallengeRequest struct {
		ID                 string  `json:"id" validate:"required,uuid"`
		ShowID             string  `json:"show_id" validate:"required,uuid"`
		Title              string  `json:"title" validate:"required,gt=0"`
		Description        string  `json:"description"`
		PrizePoolAmount    float64 `json:"prize_pool_amount" validate:"required,gt=0"`
		PlayersToStart     int32   `json:"players_to_start" validate:"required"`
		TimePerQuestionSec int64   `json:"time_per_question_sec"`
		EpisodeID          string  `json:"episode_id" validate:"uuid"`
		Kind               int32   `json:"kind"`
	}

	CheckAnswerRequest struct {
		QuestionID string `json:"question_id" validate:"required,uuid"`
		AnswerID   string `json:"answer_id" validate:"required,uuid"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetChallengeById:                   MakeGetChallengeByIdEndpoint(s),
		GetChallengesByShowId:              MakeGetChallengeByIdEndpoint(s),
		GetVerificationQuestionByEpisodeID: MakeGetVerificationQuestionByEpisodeIDEndpoint(s),
		CheckVerificationQuestionAnswer:    MakeCheckVerificationQuestionAnswerEndpoint(s, validateFunc),
		VerifyUserAccessToEpisode:          MakeVerifyUserAccessToEpisodeEndpoint(s, validateFunc),
		AddChallenge:                       MakeAddChallengeEndpoint(s, validateFunc),
		DeleteChallengeByID:                MakeDeleteChallengeByIDEndpoint(s),
		UpdateChallenge:                    MakeUpdateChallengeEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetChallengeById = mdw(e.GetChallengeById)
			e.GetChallengesByShowId = mdw(e.GetChallengesByShowId)
			e.GetVerificationQuestionByEpisodeID = mdw(e.GetVerificationQuestionByEpisodeID)
			e.CheckVerificationQuestionAnswer = mdw(e.CheckVerificationQuestionAnswer)
			e.VerifyUserAccessToEpisode = mdw(e.VerifyUserAccessToEpisode)
			e.AddChallenge = mdw(e.AddChallenge)
			e.DeleteChallengeByID = mdw(e.DeleteChallengeByID)
			e.UpdateChallenge = mdw(e.UpdateChallenge)
		}
	}

	return e
}

// MakeGetVerificationQuestionByEpisodeIDEndpoint ...
func MakeGetVerificationQuestionByEpisodeIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetVerificationQuestionByEpisodeID(ctx, id)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeCheckVerificationQuestionAnswerEndpoint ...
func MakeCheckVerificationQuestionAnswerEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CheckAnswerRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		qid, err := uuid.Parse(req.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("%w question id: %v", ErrInvalidParameter, err)
		}

		aid, err := uuid.Parse(req.AnswerID)
		if err != nil {
			return nil, fmt.Errorf("%w answer id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.CheckVerificationQuestionAnswer(ctx, qid, aid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeVerifyUserAccessToEpisodeEndpoint ...
func MakeVerifyUserAccessToEpisodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		epid, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		resp, err := s.VerifyUserAccessToEpisode(ctx, uid, epid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
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

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("could not get episode id: %w", err)
		}

		resp, err := s.AddChallenge(ctx, Challenge{
			ShowID:             showID,
			Title:              req.Title,
			Description:        req.Description,
			PrizePoolAmount:    req.PrizePoolAmount,
			Players:            req.PlayersToStart,
			TimePerQuestionSec: req.TimePerQuestionSec,
			EpisodeID:          episodeID,
			Kind:               req.Kind,
		})
		if err != nil {
			return nil, err
		}

		return resp, nil
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

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("could not get episode id: %w", err)
		}

		err = s.UpdateChallenge(ctx, Challenge{
			ID:                 id,
			ShowID:             showID,
			Title:              req.Title,
			Description:        req.Description,
			PrizePoolAmount:    req.PrizePoolAmount,
			Players:            req.PlayersToStart,
			TimePerQuestionSec: req.TimePerQuestionSec,
			EpisodeID:          episodeID,
			Kind:               req.Kind,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}
