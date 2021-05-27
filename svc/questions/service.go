package questions

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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
		CheckAnswer(ctx context.Context, arg repository.CheckAnswerParams) (sql.NullBool, error)
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
		questionID uuid.UUID `json:"question_id"`
		option     string    `json:"option"`
		isCorrect  bool      `json:"is_correct"`
	}
)

// NewService is a factory function, returns a new instance of the Service interface implementation
func NewService(qr questionsRepository) *Service {
	if qr == nil {
		log.Fatalln("question repository is not set")
	}
	return &Service{qr: qr}
}

// GetQuestionByID returns question by id
func (s *Service) GetQuestionByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	question, err := s.qr.GetQuestionByID(ctx, id)
	if !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("could not get question: %w", err)
	}

	answers, err := s.qr.GetAnswersByQuestionID(ctx, id)
	if !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("could not get question: %w", err)
	}

	return castToQuestion(question, answers), nil
}

func castToQuestion(q repository.Question, answers []repository.QuestionOption) *Question {

	question := &Question{
		ID:          q.ID,
		ChallengeID: q.ChallengeID,
		Question:    q.Question,
		Order:       q.QuestionOrder,
	}
	for i := 0; i < len(answers); i++ {
		question.QuestionOptions = append(question.QuestionOptions, QuestionOption{
			ID:         answers[i].ID,
			questionID: answers[i].QuestionID,
			option:     answers[i].QuestionOption,
			isCorrect:  answers[i].IsCorrect.Bool,
		})
	}

	return question
}

// GetQuestionByChallengeID returns questions by challenge id
func (s *Service) GetQuestionByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	questions, err := s.qr.GetQuestionsByChallengeID(ctx, id)
	if !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("could not get questions by challenge id: %w", err)
	}

	return questions, nil
}

// CheckAnswer returns question by id
func (s *Service) CheckAnswer(ctx context.Context, id, questionID uuid.UUID) (bool, error) {
	answers, err := s.qr.CheckAnswer(ctx, repository.CheckAnswerParams{
		ID:         id,
		QuestionID: questionID,
	})
	if !db.IsNotFoundError(err) {
		return false, fmt.Errorf("could not validate answer: %w", err)
	}

	return answers.Bool, nil
}
