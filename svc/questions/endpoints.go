package questions

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
		AddQuestion endpoint.Endpoint
	}

	service interface {
		AddQuestion(ctx context.Context, qw Question) (Question, error)
	}

	// AddQuestionRequest struct
	AddQuestionRequest struct {
		ID          string `json:"id" validate:"required,uuid"`
		ChallengeID string `json:"challenge_id" validate:"required,uuid"`
		Question    string `json:"question" validate:"required,gt=0"`
		Order       int32  `json:"order" validate:"required,gt=0"`
	}
)

// MakeEndpoints ...
func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		AddQuestion: MakeAddQuestionEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.AddQuestion = mdw(e.AddQuestion)
		}
	}

	return e
}

// MakeAddQuestionEndpoint ...
func MakeAddQuestionEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddQuestionRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		challengeID, err := uuid.Parse(req.ChallengeID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		resp, err := s.AddQuestion(ctx, Question{
			ChallengeID: challengeID,
			Question:    req.Question,
			Order:       req.Order,
		})
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
