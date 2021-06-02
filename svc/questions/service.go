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
		GetAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]repository.AnswerOption, error)
		GetAnswersByIDs(ctx context.Context, questionIds []uuid.UUID) ([]repository.AnswerOption, error)
		CheckAnswer(ctx context.Context, id uuid.UUID) (sql.NullBool, error)
	}

	// Question struct
	Question struct {
		ID            uuid.UUID      `json:"id"`
		ChallengeID   uuid.UUID      `json:"challenge_id"`
		Question      string         `json:"question"`
		Order         int32          `json:"order"`
		AnswerOptions []AnswerOption `json:"question_options"`
	}

	// AnswerOption struct
	AnswerOption struct {
		ID         uuid.UUID `json:"id"`
		QuestionID uuid.UUID `json:"question_id"`
		Option     string    `json:"option"`
		IsCorrect  bool      `json:"is_correct"`
	}
)

// GetQuestionByID returns question by id
func (s *Service) GetQuestionByID(ctx context.Context, id uuid.UUID) (Question, error) {
	question, err := s.qr.GetQuestionByID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return Question{}, fmt.Errorf("could not get question: %w", err)
		}
		return Question{}, fmt.Errorf("could not found question with i=%s: %w", id.String(), err)
	}

	answers, err := s.qr.GetAnswersByQuestionID(ctx, question.ID)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return Question{}, fmt.Errorf("could not get answer options for question with id=%s: %w", id.String(), err)
		}

		return Question{}, fmt.Errorf("could not found any answer options for question with id=%s: %w", id.String(), err)
	}

	return castToQuestion(question, answers), nil
}

// GetQuestionByChallengeID returns questions by challenge id
func (s *Service) GetQuestionByChallengeID(ctx context.Context, id uuid.UUID) (map[uuid.UUID]Question, error) {
	questions, err := s.qr.GetQuestionsByChallengeID(ctx, id)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get questions by challenge id: %w", err)
		}
		return nil, fmt.Errorf("could not found any questions with challenge id %s: %w", id.String(), err)
	}

	idsSlice := make([]uuid.UUID, 0, len(questions))
	for _, v := range questions {
		idsSlice = append(idsSlice, v.ID)
	}

	answers, err := s.qr.GetAnswersByIDs(ctx, idsSlice)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get answers: %w", err)
		}
		return nil, fmt.Errorf("could not found answers with ids %s: %w", id.String(), err)
	}

	m := make(map[uuid.UUID]Question, len(questions))
	for _, v := range questions {
		m[v.ID] = Question{v.ID, v.ChallengeID, v.Question,v.QuestionOrder, []AnswerOption{}}
	}

	for _, v := range answers{
		for key, val := range m{
			if v.QuestionID == key{
				val.AnswerOptions = append(val.AnswerOptions, AnswerOption{
					ID:         v.ID,
					QuestionID: v.QuestionID,
					Option:     v.AnswerOption,
					IsCorrect:  v.IsCorrect.Bool,
				})
			}
		}
	}

	return m, nil
}

// CheckAnswer checks answer
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

func castToQuestion(q repository.Question, a []repository.AnswerOption) Question {
	options := make([]AnswerOption, 0, len(a))
	for _, ao := range a {
		options = append(options, AnswerOption{
			ID:         ao.ID,
			QuestionID: ao.QuestionID,
			Option:     ao.AnswerOption,
			IsCorrect:  ao.IsCorrect.Bool,
		})
	}

	return Question{
		ID:            q.ID,
		ChallengeID:   q.ChallengeID,
		Question:      q.Question,
		Order:         q.QuestionOrder,
		AnswerOptions: options,
	}
}
