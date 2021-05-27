package questions

import (
	"context"
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
	}

	// Question struct
	Question struct {
		ID          uuid.UUID `json:"id"`
		ChallengeID uuid.UUID `json:"challenge_id"`
		Question    string    `json:"question"`
		Order       int32     `json:"order"`
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

	return question, nil
}

// GetQuestionByChallengeID returns questions by challenge id
func (s *Service) GetQuestionByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	questions, err := s.qr.GetQuestionsByChallengeID(ctx, id)
	if !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("could not get questions by challenge id: %w", err)
	}

	return questions, nil
}
