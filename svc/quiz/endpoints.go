package quiz

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetQuizLink    endpoint.Endpoint
		ParseQuizToken endpoint.Endpoint
		StoreAnswer    endpoint.Endpoint
		DeleteQuizByID endpoint.Endpoint
	}

	service interface {
		GetQuizLink(ctx context.Context, uid uuid.UUID, username string, challengeID uuid.UUID) (interface{}, error)
		ParseQuizToken(_ context.Context, token string) (*TokenPayload, error)
		StoreAnswer(ctx context.Context, userID, quizID, questionID, answerID uuid.UUID) error
		DeleteQuizByID(ctx context.Context, id uuid.UUID) error
	}

	// ConnectionURL struct
	ConnectionURL struct {
		PlayURL string `json:"play_url"`
	}

	// ParseQuizTokenRequest struct
	ParseQuizTokenRequest struct {
		Token string `json:"token" validate:"required, gt=0"`
	}

	// StoreAnswerRequest struct
	StoreAnswerRequest struct {
		UserID     string `json:"user_id" validate:"required,uuid"`
		QuizID     string `json:"quiz_id" validate:"required,uuid"`
		QuestionID string `json:"question_id" validate:"required,uuid"`
		AnswerID   string `json:"answer_id" validate:"required,uuid"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetQuizLink:    MakeGetQuizLinkEndpoint(s),
		ParseQuizToken: MakeParseQuizTokenEndpoint(s, validateFunc),
		StoreAnswer:    MakeStoreAnswerEndpoint(s, validateFunc),
		DeleteQuizByID: MakeDeleteQuizByIDEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetQuizLink = mdw(e.GetQuizLink)
			e.ParseQuizToken = mdw(e.ParseQuizToken)
			e.StoreAnswer = mdw(e.StoreAnswer)
			e.DeleteQuizByID = mdw(e.DeleteQuizByID)
		}
	}

	return e
}

// MakeGetQuizLinkEndpoint ...
func MakeGetQuizLinkEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		username, err := jwt.UsernameFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get username: %w", err)
		}

		challengeID, err := uuid.Parse(req.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		playURL, err := s.GetQuizLink(ctx, uid, username, challengeID)
		if err != nil {
			return nil, err
		}

		return ConnectionURL{
			PlayURL: fmt.Sprintf("%v", playURL),
		}, nil
	}
}

// MakeParseQuizTokenEndpoint ...
func MakeParseQuizTokenEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ParseQuizTokenRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.ParseQuizToken(ctx, req.Token)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeDeleteQuizByIDEndpoint ...
func MakeDeleteQuizByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get quiz id: %w", err)
		}

		err = s.DeleteQuizByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeStoreAnswerEndpoint ...
func MakeStoreAnswerEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreAnswerRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		aid, err := uuid.Parse(req.AnswerID)
		if err != nil {
			return nil, fmt.Errorf("could not get answer id: %w", err)
		}

		qzid, err := uuid.Parse(req.QuizID)
		if err != nil {
			return nil, fmt.Errorf("could not get quiz id: %w", err)
		}

		qid, err := uuid.Parse(req.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("could not get question id: %w", err)
		}

		uid, err := uuid.Parse(req.UserID)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		err = s.StoreAnswer(ctx, uid, qzid, qid, aid)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}
