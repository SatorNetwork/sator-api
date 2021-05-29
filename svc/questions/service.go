package questions

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/questions/repository"
	"github.com/google/uuid"
)

type (
	//Service struct
	Service struct {
		qr questionsRepository
	}

	questionsRepository interface {
		AddQuestion(ctx context.Context, arg repository.AddQuestionParams) (repository.Question, error)
		GetQuestionByID(ctx context.Context, id uuid.UUID) (repository.Question, error)
		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) ([]repository.Question, error)
		GetAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]repository.QuestionOption, error)
		CheckAnswer(ctx context.Context, id uuid.UUID) (sql.NullBool, error)
	}

	// Question struct
	Question struct {
		ID              uuid.UUID        `json:"id"`
		ChallengeID     uuid.UUID        `json:"challenge_id"`
		Question        string           `json:"question"`
		Order           int32            `json:"order"`
		QuestionOptions []QuestionOption `json:"question_options"`
	}

	// QuestionOption struct
	QuestionOption struct {
		ID         uuid.UUID `json:"id"`
		QuestionID uuid.UUID `json:"question_id"`
		Option     string    `json:"option"`
		IsCorrect  bool      `json:"is_correct"`
	}
)

// GetQuestionByID returns question by id
func (s *Service) GetQuestionByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	question, err := s.qr.GetQuestionByID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get question: %w", err)
		}
		return nil, fmt.Errorf("could not found question with i=%s: %w", id.String(), err)
	}

	answers, err := s.qr.GetAnswersByQuestionID(ctx, question.ID)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get answer options for question with id=%s: %w", id.String(), err)
		}

		return nil, fmt.Errorf("could not found any answer options for question with id=%s: %w", id.String(), err)
	}

	return castToQuestion(question, answers), nil
}

// GetQuestionByChallengeID returns questions by challenge id
// TODO: needs refactoring! Make query that will return required slice.
func (s *Service) GetQuestionByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	questions, err := s.qr.GetQuestionsByChallengeID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get questions by challenge id: %w", err)
		}
		return nil, fmt.Errorf("could not found any questions with challenge id %s: %w", id.String(), err)
	}

	result := make([]Question, 0, len(questions))
	for _, q := range questions {
		answers, err := s.qr.GetAnswersByQuestionID(ctx, q.ID)
		if err != nil {
			if !db.IsNotFoundError(err) {
				return nil, fmt.Errorf("could not get answer: %w", err)
			}
			return nil, fmt.Errorf("could not found answer with id %s: %w", id.String(), err)
		}

		question := castToQuestion(q, answers)
		result = append(result, question)
	}

	return result, nil
}

// CheckAnswer returns question by id
func (s *Service) CheckAnswer(ctx context.Context, id uuid.UUID) (bool, error) {
	answers, err := s.qr.CheckAnswer(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return false, fmt.Errorf("could not validate answer: %w", err)
		}
		return false, fmt.Errorf("could not found question with id %s: %w", id, err)
	}

	return answers.Bool, nil
}

func castToQuestion(q repository.Question, a []repository.QuestionOption) Question {
	options := make([]QuestionOption, 0, len(a))
	for _, ao := range a {
		options = append(options, QuestionOption{
			ID:         ao.ID,
			QuestionID: ao.QuestionID,
			Option:     ao.QuestionOption,
			IsCorrect:  ao.IsCorrect.Bool,
		})
	}

	return Question{
		ID:              q.ID,
		ChallengeID:     q.ChallengeID,
		Question:        q.Question,
		Order:           q.QuestionOrder,
		QuestionOptions: options,
	}
}
