package quiz

import (
	"context"

	"github.com/SatorNetwork/sator-api/svc/quiz/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		repo quizRepository
	}

	quizRepository interface {
		AddNewChallegeRoom(ctx context.Context, arg repository.AddNewChallegeRoomParams) (repository.ChallengeRoom, error)
		GetChallengeRoomByID(ctx context.Context, id uuid.UUID) (repository.ChallengeRoom, error)
		UpdateChallengeRoomStatus(ctx context.Context, arg repository.UpdateChallengeRoomStatusParams) error
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo quizRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// Start function registers user in waiting room
func (s *Service) Start(ctx context.Context, uid, challengeID uuid.UUID) (interface{}, error) {
	return nil, nil
}

// Play function registers user in quiz hub
func (s *Service) Play(ctx context.Context, uid, roomID uuid.UUID) (interface{}, error) {
	return nil, nil
}
