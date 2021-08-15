package challenge

import (
	"context"
	"fmt"
	"strconv"

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
		GetChallengeByID(ctx context.Context, id uuid.UUID) (interface{}, error)
		GetChallengesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error)
		AddChallenge(ctx context.Context, ch Challenge) (Challenge, error)
		DeleteChallengeByID(ctx context.Context, id uuid.UUID) error
		UpdateChallenge(ctx context.Context, ch Challenge) error

		// FIXME: needs refactoring!
		GetVerificationQuestionByEpisodeID(ctx context.Context, episodeID, userID uuid.UUID) (interface{}, error)
		CheckVerificationQuestionAnswer(ctx context.Context, questionID, answerID, userID uuid.UUID) (interface{}, error)
		VerifyUserAccessToEpisode(ctx context.Context, uid, eid uuid.UUID) (interface{}, error)

		AddQuestion(ctx context.Context, qw Question) (Question, error)
		DeleteQuestionByID(ctx context.Context, id uuid.UUID) error
		GetQuestionByID(ctx context.Context, id uuid.UUID) (Question, error)
		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error)
		UpdateQuestion(ctx context.Context, qw Question) error

		AddQuestionOption(ctx context.Context, ao AnswerOption) (AnswerOption, error)
		DeleteAnswerByID(ctx context.Context, id, questionID uuid.UUID) error
		CheckAnswer(ctx context.Context, aid, qid uuid.UUID) (bool, error)
		UpdateAnswer(ctx context.Context, ao AnswerOption) error
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

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetChallengeById:                   MakeGetChallengeByIdEndpoint(s),
		GetChallengesByShowId:              MakeGetChallengeByIdEndpoint(s),
		GetVerificationQuestionByEpisodeID: MakeGetVerificationQuestionByEpisodeIDEndpoint(s),
		CheckVerificationQuestionAnswer:    MakeCheckVerificationQuestionAnswerEndpoint(s, validateFunc),
		VerifyUserAccessToEpisode:          MakeVerifyUserAccessToEpisodeEndpoint(s),
		AddChallenge:                       MakeAddChallengeEndpoint(s, validateFunc),
		DeleteChallengeByID:                MakeDeleteChallengeByIDEndpoint(s),
		UpdateChallenge:                    MakeUpdateChallengeEndpoint(s, validateFunc),

		AddQuestion:               MakeAddQuestionEndpoint(s, validateFunc),
		DeleteQuestionByID:        MakeDeleteQuestionByIDEndpoint(s),
		UpdateQuestion:            MakeUpdateQuestionEndpoint(s, validateFunc),
		GetQuestionByID:           MakeGetQuestionByIDEndpoint(s),
		GetQuestionsByChallengeID: MakeGetQuestionsByChallengeIDEndpoint(s),

		AddQuestionOption: MakeAddQuestionOptionEndpoint(s, validateFunc),
		DeleteAnswerByID:  MakeDeleteAnswerByIDEndpoint(s, validateFunc),
		CheckAnswer:       MakeCheckAnswerEndpoint(s, validateFunc),
		UpdateAnswer:      MakeUpdateAnswerEndpoint(s, validateFunc),
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

// MakeGetVerificationQuestionByEpisodeIDEndpoint ...
func MakeGetVerificationQuestionByEpisodeIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetVerificationQuestionByEpisodeID(ctx, id, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeCheckVerificationQuestionAnswerEndpoint ...
func MakeCheckVerificationQuestionAnswerEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

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

		resp, err := s.CheckVerificationQuestionAnswer(ctx, qid, aid, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeVerifyUserAccessToEpisodeEndpoint ...
func MakeVerifyUserAccessToEpisodeEndpoint(s service) endpoint.Endpoint {
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
func MakeCheckAnswerEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CheckAnswerRequest)
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

		resp, err := s.CheckAnswer(ctx, answerID, questionID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
