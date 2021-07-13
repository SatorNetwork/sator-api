package questions

import (
	"context"
	"fmt"
	"strconv"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		AddQuestion               endpoint.Endpoint
		DeleteQuestionByID        endpoint.Endpoint
		GetQuestionByID           endpoint.Endpoint
		GetQuestionsByChallengeID endpoint.Endpoint
		UpdateQuestion            endpoint.Endpoint

		AddQuestionOption endpoint.Endpoint
		DeleteAnswerByID  endpoint.Endpoint
		CheckAnswer       endpoint.Endpoint
		UpdateAnswer      endpoint.Endpoint
	}

	service interface {
		AddQuestion(ctx context.Context, qw Question) (Question, error)
		DeleteQuestionByID(ctx context.Context, id uuid.UUID) error
		GetQuestionByID(ctx context.Context, id uuid.UUID) (Question, error)
		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error)
		UpdateQuestion(ctx context.Context, qw Question) error

		AddQuestionOption(ctx context.Context, ao AnswerOption) (AnswerOption, error)
		DeleteAnswerByID(ctx context.Context, id, questionID uuid.UUID) error
		CheckAnswer(ctx context.Context, id uuid.UUID) (bool, error)
		UpdateAnswer(ctx context.Context, ao AnswerOption) error
	}

	// AddQuestionRequest struct
	AddQuestionRequest struct {
		ChallengeID string `json:"challenge_id" validate:"required,uuid"`
		Question    string `json:"question" validate:"required,gt=0"`
		Order       int32  `json:"order" validate:"required,gt=0"`
	}

	// UpdateQuestionRequest struct
	UpdateQuestionRequest struct {
		ID          string `json:"id" validate:"required,uuid"`
		ChallengeID string `json:"challenge_id" validate:"required,uuid"`
		Question    string `json:"question" validate:"required,gt=0"`
		Order       int32  `json:"order" validate:"required,gt=0"`
	}

	// AnswerOptionRequest struct
	AnswerOptionRequest struct {
		QuestionID string `json:"question_id" validate:"required,uuid"`
		Option     string `json:"option" validate:"required,gt=0"`
		IsCorrect  string `json:"is_correct" validate:"required,gt=0"`
	}

	// UpdateAnswerRequest struct
	UpdateAnswerRequest struct {
		ID         string `json:"id" validate:"required,uuid"`
		QuestionID string `json:"question_id" validate:"required,uuid"`
		Option     string `json:"option" validate:"required,gt=0"`
		IsCorrect  string `json:"is_correct" validate:"required,gt=0"`
	}

	// DeleteAnswerByIDRequest struct
	DeleteAnswerByIDRequest struct {
		AnswerID   string `json:"answer_id"`
		QuestionID string `json:"question_id"`
	}
)

// MakeEndpoints ...
func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		AddQuestion:               MakeAddQuestionEndpoint(s, validateFunc),
		DeleteQuestionByID:        MakeDeleteQuestionByIDEndpoint(s),
		UpdateQuestion:            MakeUpdateQuestionEndpoint(s, validateFunc),
		GetQuestionByID:           MakeGetQuestionByIDEndpoint(s),
		GetQuestionsByChallengeID: MakeGetQuestionsByChallengeIDEndpoint(s),

		AddQuestionOption: MakeAddQuestionOptionEndpoint(s, validateFunc),
		DeleteAnswerByID:  MakeDeleteAnswerByIDEndpoint(s, validateFunc),
		CheckAnswer:       MakeCheckAnswerEndpoint(s),
		UpdateAnswer:      MakeUpdateAnswerEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.AddQuestion = mdw(e.AddQuestion)
			e.DeleteQuestionByID = mdw(e.DeleteQuestionByID)
			e.UpdateQuestion = mdw(e.UpdateQuestion)
			e.GetQuestionByID = mdw(e.GetQuestionByID)
			e.GetQuestionsByChallengeID = mdw(e.GetQuestionsByChallengeID)

			e.AddQuestionOption = mdw(e.AddQuestionOption)
			e.DeleteAnswerByID = mdw(e.DeleteAnswerByID)
			e.CheckAnswer = mdw(e.CheckAnswer)
			e.UpdateAnswer = mdw(e.UpdateAnswer)
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

// MakeAddQuestionOptionEndpoint ...
func MakeAddQuestionOptionEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AnswerOptionRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		questionID, err := uuid.Parse(req.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		isCorrect, err := strconv.ParseBool(req.IsCorrect)
		if err != nil {
			return nil, fmt.Errorf("could not parse bool from string %w", err)
		}

		resp, err := s.AddQuestionOption(ctx, AnswerOption{
			QuestionID: questionID,
			Option:     req.Option,
			IsCorrect:  isCorrect,
		})
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeDeleteQuestionByIDEndpoint ...
func MakeDeleteQuestionByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get question id: %w", err)
		}

		err = s.DeleteQuestionByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%w question id: %v", ErrInvalidParameter, err)
		}

		return true, nil
	}
}

// MakeDeleteAnswerByIDEndpoint ...
func MakeDeleteAnswerByIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteAnswerByIDRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		answerID, err := uuid.Parse(req.AnswerID)
		if err != nil {
			return nil, fmt.Errorf("could not get answer id: %w", err)
		}

		questionID, err := uuid.Parse(req.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("could not get question id: %w", err)
		}

		err = s.DeleteAnswerByID(ctx, answerID, questionID)
		if err != nil {
			return nil, fmt.Errorf("%w answer id: %v", ErrInvalidParameter, err)
		}

		return true, nil
	}
}

// MakeUpdateQuestionEndpoint ...
func MakeUpdateQuestionEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateQuestionRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get question id: %w", err)
		}

		challengeID, err := uuid.Parse(req.ChallengeID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		if err = s.UpdateQuestion(ctx, Question{
			ID:          id,
			ChallengeID: challengeID,
			Question:    req.Question,
			Order:       req.Order,
		}); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeUpdateAnswerEndpoint ...
func MakeUpdateAnswerEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateAnswerRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get answer id: %w", err)
		}

		questionID, err := uuid.Parse(req.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		isCorrect, err := strconv.ParseBool(req.IsCorrect)
		if err != nil {
			return nil, fmt.Errorf("could not parse bool from string %w", err)
		}

		if err = s.UpdateAnswer(ctx, AnswerOption{
			ID:         id,
			QuestionID: questionID,
			Option:     req.Option,
			IsCorrect:  isCorrect,
		}); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeGetQuestionByIDEndpoint ...
func MakeGetQuestionByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		questionID, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get question id: %w", err)
		}

		resp, err := s.GetQuestionByID(ctx, questionID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetQuestionsByChallengeIDEndpoint ...
func MakeGetQuestionsByChallengeIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		challengeID, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		resp, err := s.GetQuestionsByChallengeID(ctx, challengeID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeCheckAnswerEndpoint ...
func MakeCheckAnswerEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		answerID, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get answer id: %w", err)
		}

		resp, err := s.CheckAnswer(ctx, answerID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
