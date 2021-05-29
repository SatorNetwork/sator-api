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

		return nil, fmt.Errorf("question with id %w not found", id)
	}

	answers, err := s.qr.GetAnswersByQuestionID(ctx, question.ID)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get answer: %w", err)
		}

		return nil, fmt.Errorf("answer with id %w not found", id)
	}

	return castToQuestion(question, answers), nil
}

func castToQuestion(q repository.Question, a []repository.QuestionOption) *Question {
	question := &Question{
		ID:          q.ID,
		ChallengeID: q.ChallengeID,
		Question:    q.Question,
		Order:       q.QuestionOrder,
	}
	for i := 0; i < len(a); i++ {
		question.QuestionOptions = append(question.QuestionOptions, QuestionOption{
			ID:         a[i].ID,
			QuestionID: a[i].QuestionID,
			Option:     a[i].QuestionOption,
			IsCorrect:  a[i].IsCorrect.Bool,
		})
	}

	return question
}

// GetQuestionByChallengeID returns questions by challenge id
func (s *Service) GetQuestionByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error) {  // TODO: Make query that will return required slice.
	var result []Question
	questions, err := s.qr.GetQuestionsByChallengeID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get questions by challenge id: %w", err)
		}

		return nil, fmt.Errorf("questions with challenge id %w not found", id)
	}

	for _, q := range questions {
		answers, err := s.qr.GetAnswersByQuestionID(ctx, q.ID)
		if err != nil {
			if !db.IsNotFoundError(err) {
				return nil, fmt.Errorf("could not get answer: %w", err)
			}

			return nil, fmt.Errorf("answer with id %w not found", id)
		}

		question := castToQuestion(q, answers)
		result = append(result, *question)
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

		return false, fmt.Errorf("question with id %w not found", id)
	}

	return answers.Bool, nil
}
