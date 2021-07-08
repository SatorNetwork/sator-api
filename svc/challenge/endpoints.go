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

		GetVerificationQuestionByEpisodeID endpoint.Endpoint
		CheckVerificationQuestionAnswer    endpoint.Endpoint
		VerifyUserAccessToEpisode          endpoint.Endpoint
	}

	service interface {
		GetChallengeByID(ctx context.Context, id uuid.UUID) (interface{}, error)
		GetChallengesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error)
		// FIXME: needs refactoring!
		GetVerificationQuestionByEpisodeID(ctx context.Context, episodeID uuid.UUID) (interface{}, error)
		CheckVerificationQuestionAnswer(ctx context.Context, qid, aid uuid.UUID) (interface{}, error)
		VerifyUserAccessToEpisode(ctx context.Context, uid, eid uuid.UUID) (interface{}, error)
	}

	GetChallengeByIdRequest struct {
		ID string `json:"id" validate:"required,uuid"`
	}

	CheckAnswerRequest struct {
		QuestionID string `json:"question_id" validate:"required,uuid"`
		AnswerID   string `json:"answer_id" validate:"required,uuid"`
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

// MakeGetVerificationQuestionByEpisodeIDEndpoint ...
func MakeGetVerificationQuestionByEpisodeIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
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

		qid, err := uuid.Parse(req.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("%w question id: %v", ErrInvalidParameter, err)
		}

		aid, err := uuid.Parse(req.QuestionID)
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
