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
		repo      rewardsRepository
		assetName string
	}

	rewardsRepository interface {
		AddReward(ctx context.Context, arg repository.AddRewardParams) error
	}

	ClaimRewardsResult struct {
		DisplayAmount   string  `json:"amount"`
		TransactionURL  string  `json:"transaction_url"`
		Amount          float64 `json:"-"`
		TransactionHash string  `json:"-"`
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo rewardsRepository) *Service {
	return &Service{
		repo:      repo,
		assetName: "SAO",
	}
}

// AddReward ...
func (s *Service) AddReward(ctx context.Context, uid uuid.UUID, amount float64, qid uuid.UUID) error {
	if err := s.repo.AddReward(ctx, repository.AddRewardParams{
		UserID: uid,
		QuizID: qid,
		Amount: amount,
	}); err != nil {
		return fmt.Errorf("could not add reward for user_id=%s and quiz_id=%s: %w", uid.String(), qid.String(), err)
	}
	return nil
}

// ClaimRewards ...
func (s *Service) ClaimRewards(ctx context.Context, uid uuid.UUID) (ClaimRewardsResult, error) {
	return ClaimRewardsResult{
		Amount:        83.54,
		DisplayAmount: fmt.Sprintf("%.2f %s", 83.54, s.assetName),
	}, nil
}
