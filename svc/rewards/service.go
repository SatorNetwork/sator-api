package rewards

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/svc/rewards/repository"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

type (
	// Service struct
	Service struct {
		repo      rewardsRepository
		assetName string
		ws        walletService
	}

	Winner struct {
		UserID uuid.UUID
		Points int
	}

	rewardsRepository interface {
		AddReward(ctx context.Context, arg repository.AddRewardParams) error
		GetUnWithdrawnRewards(ctx context.Context, arg repository.GetUnWithdrawnRewardsParams) ([]repository.Reward, error)
		Withdraw(ctx context.Context, userID uuid.UUID) error
	}

	ClaimRewardsResult struct {
		DisplayAmount   string  `json:"amount"`
		TransactionURL  string  `json:"transaction_url"`
		Amount          float64 `json:"-"`
		TransactionHash string  `json:"-"`
	}

	walletService interface {
		SendToWallet(ctx context.Context, userID uuid.UUID, amount float64) (string, error)
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

// ClaimRewards send rewards to user by it and sets them to withdrawn.
func (s *Service) ClaimRewards(ctx context.Context, uid uuid.UUID) (ClaimRewardsResult, error) {
	rewards, err := s.repo.GetUnWithdrawnRewards(ctx, repository.GetUnWithdrawnRewardsParams{
		UserID:    uuid.UUID{},
		Withdrawn: false,
	})
	if err != nil {
		return ClaimRewardsResult{}, err
	}
	var amount float64

	for _, reward := range rewards {
		amount += reward.Amount
	}

	txHash, err := s.ws.SendToWallet(ctx, uid, amount)
	if err != nil {
		return ClaimRewardsResult{}, err
	}

	err = s.repo.Withdraw(ctx, uid)
	if err != nil {
		return ClaimRewardsResult{}, err
	}

	return ClaimRewardsResult{
		Amount:          amount,
		DisplayAmount:   fmt.Sprintf("%.2f %s", amount, s.assetName),
		TransactionHash: txHash,
	}, nil
}

// DistributeRewards split rewards among users, store into db.
func (s *Service) DistributeRewards(ctx context.Context, prizePool float64, winners []Winner, qid uuid.UUID) (err error) {
	var totalPoints int

	for _, winner := range winners {
		totalPoints += winner.Points
	}
	pointCost := prizePool / float64(totalPoints)

	for _, winner := range winners {
		err = errs.Combine(err, s.AddReward(ctx, winner.UserID, pointCost*float64(winner.Points), qid))
	}

	return err
}
