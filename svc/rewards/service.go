package rewards

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/svc/rewards/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		repo rewardsRepository
	}

	rewardsRepository interface {
		AddReward(ctx context.Context, arg repository.AddRewardParams) error
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo rewardsRepository) *Service {
	return &Service{repo: repo}
}

// AddReward ...
func (s *Service) AddReward(ctx context.Context, uid uuid.UUID, amount float64, qid uuid.UUID) error {
	if err := s.repo.AddReward(ctx, repository.AddRewardParams{}); err != nil {
		return fmt.Errorf("could not add reward for user_id=%s and quiz_id=%s: %w", uid.String(), qid.String(), err)
	}
	return nil
}
