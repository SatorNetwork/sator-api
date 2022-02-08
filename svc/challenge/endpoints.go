package challenge

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/rbac"
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

		GetVerificationQuestionByEpisodeID     endpoint.Endpoint
		CheckVerificationQuestionAnswer        endpoint.Endpoint
		VerifyUserAccessToEpisode              endpoint.Endpoint
		GetAttemptsLeftForVerificationQuestion endpoint.Endpoint

		AddQuestion               endpoint.Endpoint
		DeleteQuestionByID        endpoint.Endpoint
		GetQuestionByID           endpoint.Endpoint
		GetQuestionsByChallengeID endpoint.Endpoint
		UpdateQuestion            endpoint.Endpoint

		AddQuestionOption endpoint.Endpoint
		DeleteAnswerByID  endpoint.Endpoint
		CheckAnswer       endpoint.Endpoint
		UpdateAnswer      endpoint.Endpoint

		UnlockEpisode endpoint.Endpoint
	}

	service interface {
		GetChallengeByID(ctx context.Context, challengeID, userID uuid.UUID) (interface{}, error)
		GetChallengesByShowID(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error)
		AddChallenge(ctx context.Context, ch Challenge) (Challenge, error)
		DeleteChallengeByID(ctx context.Context, id uuid.UUID) error
		UpdateChallenge(ctx context.Context, ch Challenge) error

		GetVerificationQuestionByEpisodeID(ctx context.Context, episodeID, userID uuid.UUID) (interface{}, error)
		CheckVerificationQuestionAnswer(ctx context.Context, questionID, answerID, userID uuid.UUID) (interface{}, error)
		VerifyUserAccessToEpisode(ctx context.Context, uid, eid uuid.UUID) (interface{}, error)
		GetAttemptsLeftForVerificationQuestion(ctx context.Context, episodeID, userID uuid.UUID) (int64, error)

		AddQuestion(ctx context.Context, qw Question) (Question, error)
		DeleteQuestionByID(ctx context.Context, id uuid.UUID) error
		GetQuestionByID(ctx context.Context, id uuid.UUID) (Question, error)
		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error)
		UpdateQuestion(ctx context.Context, qw Question) error

		AddQuestionOption(ctx context.Context, ao AnswerOption) (AnswerOption, error)
		DeleteAnswerByID(ctx context.Context, id, questionID uuid.UUID) error
		DeleteAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) error
		CheckAnswer(ctx context.Context, aid, qid uuid.UUID) (bool, error)
		UpdateAnswer(ctx context.Context, ao AnswerOption) error

		UnlockEpisode(ctx context.Context, uid, episodeID uuid.UUID, unlockOption string) error
	}

	// AddChallengeRequest struct
	AddChallengeRequest struct {
		ShowID             string  `json:"show_id" validate:"required,uuid"`
		Title              string  `json:"title" validate:"required,gt=0"`
		Description        string  `json:"description"`
		PrizePoolAmount    float64 `json:"prize_pool_amount" validate:"gte=0"`
		PlayersToStart     int     `json:"players_to_start" validate:"required,gt=0"`
		TimePerQuestionSec int     `json:"time_per_question_sec" validate:"required,gt=0"`
		EpisodeID          string  `json:"episode_id,omitempty"`
		Kind               int     `json:"kind"`
		UserMaxAttempts    int     `json:"user_max_attempts" validate:"required,gt=0"`
		MaxWinners         int     `json:"max_winners"`
		QuestionsPerGame   int     `json:"questions_per_game"`
		MinCorrectAnswers  int     `json:"min_correct_answers"`
	}

	// UpdateChallengeRequest struct
	UpdateChallengeRequest struct {
		ID                 string  `json:"id" validate:"required,uuid"`
		ShowID             string  `json:"show_id" validate:"required,uuid"`
		Title              string  `json:"title" validate:"required,gt=0"`
		Description        string  `json:"description"`
		PrizePoolAmount    float64 `json:"prize_pool_amount" validate:"gte=0"`
		PlayersToStart     int     `json:"players_to_start" validate:"required,gt=0"`
		TimePerQuestionSec int     `json:"time_per_question_sec" validate:"required,gt=0"`
		EpisodeID          string  `json:"episode_id"`
		Kind               int     `json:"kind"`
		UserMaxAttempts    int     `json:"user_max_attempts" validate:"required,gt=0"`
		MaxWinners         int     `json:"max_winners"`
		QuestionsPerGame   int     `json:"questions_per_game"`
		MinCorrectAnswers  int     `json:"min_correct_answers"`
	}

	CheckAnswerRequest struct {
		QuestionID string `json:"question_id" validate:"required,uuid"`
		AnswerID   string `json:"answer_id" validate:"required,uuid"`
	}

	// AddQuestionRequest struct
	AddQuestionRequest struct {
		ChallengeID   string                `json:"challenge_id" validate:"required,uuid"`
		Question      string                `json:"question" validate:"required,gt=0"`
		Order         int32                 `json:"order,omitempty"`
		AnswerOptions []AnswerOptionRequest `json:"answer_options,omitempty"`
	}

	// UpdateQuestionRequest struct
	UpdateQuestionRequest struct {
		ID            string                `json:"id" validate:"required,uuid"`
		ChallengeID   string                `json:"challenge_id" validate:"required,uuid"`
		Question      string                `json:"question" validate:"required,gt=0"`
		Order         int32                 `json:"order,omitempty"`
		AnswerOptions []AnswerOptionRequest `json:"answer_options,omitempty"`
	}

	// AnswerOptionRequest struct
	AnswerOptionRequest struct {
		QuestionID string `json:"question_id,omitempty"`
		Option     string `json:"option" validate:"required,gt=0"`
		IsCorrect  bool   `json:"is_correct" validate:"required"`
	}

	// UpdateAnswerRequest struct
	UpdateAnswerRequest struct {
		ID         string `json:"id" validate:"required,uuid"`
		QuestionID string `json:"question_id" validate:"required,uuid"`
		Option     string `json:"option" validate:"required,gt=0"`
		IsCorrect  bool   `json:"is_correct" validate:"required"`
	}

	// DeleteAnswerByIDRequest struct
	DeleteAnswerByIDRequest struct {
		AnswerID   string `json:"answer_id"`
		QuestionID string `json:"question_id"`
	}

	// UnlockEpisodeRequest ...
	UnlockEpisodeRequest struct {
		EpisodeID string `json:"episode_id" validate:"required,uuid"`
		Option    string `json:"option" validate:"required"`
	}

	//	GetAttemptsLeftForVerificationQuestionResponse...
	GetAttemptsLeftForVerificationQuestionResponse struct {
		AttemptsLeft int64 `json:"attempts_left"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetChallengeById:                       MakeGetChallengeByIdEndpoint(s),
		GetChallengesByShowId:                  MakeGetChallengeByIdEndpoint(s),
		GetVerificationQuestionByEpisodeID:     MakeGetVerificationQuestionByEpisodeIDEndpoint(s),
		CheckVerificationQuestionAnswer:        MakeCheckVerificationQuestionAnswerEndpoint(s, validateFunc),
		VerifyUserAccessToEpisode:              MakeVerifyUserAccessToEpisodeEndpoint(s),
		AddChallenge:                           MakeAddChallengeEndpoint(s, validateFunc),
		DeleteChallengeByID:                    MakeDeleteChallengeByIDEndpoint(s),
		UpdateChallenge:                        MakeUpdateChallengeEndpoint(s, validateFunc),
		GetAttemptsLeftForVerificationQuestion: MakeGetAttemptsLeftForVerificationQuestionEndpoint(s),

		AddQuestion:               MakeAddQuestionEndpoint(s, validateFunc),
		DeleteQuestionByID:        MakeDeleteQuestionByIDEndpoint(s),
		UpdateQuestion:            MakeUpdateQuestionEndpoint(s, validateFunc),
		GetQuestionByID:           MakeGetQuestionByIDEndpoint(s),
		GetQuestionsByChallengeID: MakeGetQuestionsByChallengeIDEndpoint(s),

		AddQuestionOption: MakeAddQuestionOptionEndpoint(s, validateFunc),
		DeleteAnswerByID:  MakeDeleteAnswerByIDEndpoint(s, validateFunc),
		CheckAnswer:       MakeCheckAnswerEndpoint(s, validateFunc),
		UpdateAnswer:      MakeUpdateAnswerEndpoint(s, validateFunc),

		UnlockEpisode: MakeUnlockEpisodeEndpoint(s, validateFunc),
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
			e.GetAttemptsLeftForVerificationQuestion = mdw(e.GetAttemptsLeftForVerificationQuestion)

			e.AddQuestion = mdw(e.AddQuestion)
			e.DeleteQuestionByID = mdw(e.DeleteQuestionByID)
			e.UpdateQuestion = mdw(e.UpdateQuestion)
			e.GetQuestionByID = mdw(e.GetQuestionByID)
			e.GetQuestionsByChallengeID = mdw(e.GetQuestionsByChallengeID)

			e.AddQuestionOption = mdw(e.AddQuestionOption)
			e.DeleteAnswerByID = mdw(e.DeleteAnswerByID)
			e.CheckAnswer = mdw(e.CheckAnswer)
			e.UpdateAnswer = mdw(e.UpdateAnswer)

			e.UnlockEpisode = mdw(e.UnlockEpisode)
		}
	}

	return e
}

// MakeGetVerificationQuestionByEpisodeIDEndpoint ...
func MakeGetVerificationQuestionByEpisodeIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w episode id: %v", ErrInvalidParameter, err)
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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetChallengeByID(ctx, id, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeAddChallengeEndpoint ...
func MakeAddChallengeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(AddChallengeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		payload := Challenge{
			ShowID:             showID,
			Title:              req.Title,
			Description:        req.Description,
			PrizePoolAmount:    req.PrizePoolAmount,
			Players:            int32(req.PlayersToStart),
			TimePerQuestionSec: int32(req.TimePerQuestionSec),
			Kind:               int32(req.Kind),
			UserMaxAttempts:    int32(req.UserMaxAttempts),
			MaxWinners:         int32(req.MaxWinners),
			QuestionsPerGame:   int32(req.QuestionsPerGame),
			MinCorrectAnswers:  int32(req.MinCorrectAnswers),
		}

		if req.EpisodeID != "" && req.EpisodeID != uuid.Nil.String() {
			episodeID, err := uuid.Parse(req.EpisodeID)
			if err != nil {
				return nil, fmt.Errorf("could not get episode id: %w", err)
			}
			payload.EpisodeID = &episodeID
		}

		resp, err := s.AddChallenge(ctx, payload)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeDeleteChallengeByIDEndpoint ...
func MakeDeleteChallengeByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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

		payload := Challenge{
			ID:                 id,
			ShowID:             showID,
			Title:              req.Title,
			Description:        req.Description,
			PrizePoolAmount:    req.PrizePoolAmount,
			Players:            int32(req.PlayersToStart),
			TimePerQuestionSec: int32(req.TimePerQuestionSec),
			Kind:               int32(req.Kind),
			UserMaxAttempts:    int32(req.UserMaxAttempts),
			MaxWinners:         int32(req.MaxWinners),
			QuestionsPerGame:   int32(req.QuestionsPerGame),
			MinCorrectAnswers:  int32(req.MinCorrectAnswers),
		}

		if req.EpisodeID != "" && req.EpisodeID != uuid.Nil.String() {
			episodeID, err := uuid.Parse(req.EpisodeID)
			if err != nil {
				return nil, fmt.Errorf("could not get episode id: %w", err)
			}
			payload.EpisodeID = &episodeID
		}

		err = s.UpdateChallenge(ctx, payload)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeAddQuestionEndpoint ...
func MakeAddQuestionEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(AddQuestionRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		challengeID, err := uuid.Parse(req.ChallengeID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		if n := len(req.AnswerOptions); n >= 2 && n <= 4 {
			resp, err := s.AddQuestion(ctx, Question{
				ChallengeID: challengeID,
				Question:    req.Question,
				Order:       req.Order,
			})
			if err != nil {
				return nil, err
			}

			answerOptions := make([]AnswerOption, 0, len(req.AnswerOptions))

			for _, answ := range req.AnswerOptions {
				ao, err := s.AddQuestionOption(ctx, AnswerOption{
					QuestionID: resp.ID,
					Option:     answ.Option,
					IsCorrect:  answ.IsCorrect,
				})
				if err != nil {
					return nil, err
				}

				answerOptions = append(answerOptions, ao)
			}

			resp.AnswerOptions = answerOptions

			return resp, nil
		}

		return nil, fmt.Errorf("number of answer options must be from 2 to 4")
	}
}

// MakeAddQuestionOptionEndpoint ...
func MakeAddQuestionOptionEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(AnswerOptionRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		questionID, err := uuid.Parse(req.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge id: %w", err)
		}

		resp, err := s.AddQuestionOption(ctx, AnswerOption{
			QuestionID: questionID,
			Option:     req.Option,
			IsCorrect:  req.IsCorrect,
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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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

		if n := len(req.AnswerOptions); n >= 2 && n <= 4 {
			if err = s.UpdateQuestion(ctx, Question{
				ID:          id,
				ChallengeID: challengeID,
				Question:    req.Question,
				Order:       req.Order,
			}); err != nil {
				return nil, err
			}

			if err := s.DeleteAnswersByQuestionID(ctx, id); err != nil {
				return nil, err
			}

			for _, answ := range req.AnswerOptions {
				if _, err := s.AddQuestionOption(ctx, AnswerOption{
					QuestionID: id,
					Option:     answ.Option,
					IsCorrect:  answ.IsCorrect,
				}); err != nil {
					return nil, err
				}
			}
			return true, nil
		}

		return nil, fmt.Errorf("number of answer options must be from 2 to 4")
	}
}

// MakeUpdateAnswerEndpoint ...
func MakeUpdateAnswerEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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

		if err = s.UpdateAnswer(ctx, AnswerOption{
			ID:         id,
			QuestionID: questionID,
			Option:     req.Option,
			IsCorrect:  req.IsCorrect,
		}); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeGetQuestionByIDEndpoint ...
func MakeGetQuestionByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		// FIXME: is allowed roles correct???
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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

// MakeUnlockEpisodeEndpoint ...
// TODO: remove it, added for demo
func MakeUnlockEpisodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// FIXME: is allowed roles correct???
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(UnlockEpisodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("could not get episode id: %w", err)
		}

		if err := s.UnlockEpisode(ctx, uid, episodeID, req.Option); err != nil {
			return false, err
		}

		resp, err := s.VerifyUserAccessToEpisode(ctx, uid, episodeID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetAttemptsLeftForVerificationQuestionEndpoint ...
func MakeGetAttemptsLeftForVerificationQuestionEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w episode id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetAttemptsLeftForVerificationQuestion(ctx, id, uid)
		if err != nil {
			return nil, err
		}

		return GetAttemptsLeftForVerificationQuestionResponse{AttemptsLeft: resp}, nil
	}
}
