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
		ws        walletService
	}

	Winner struct {
		UserID uuid.UUID
		Points int
	}

	rewardsRepository interface {
		AddReward(ctx context.Context, arg repository.AddRewardParams) error
		Withdraw(ctx context.Context, userID uuid.UUID) error
		GetTotalAmount(ctx context.Context, userID uuid.UUID) (float64, error)
	}

	ClaimRewardsResult struct {
		DisplayAmount   string  `json:"amount"`
		TransactionURL  string  `json:"transaction_url"`
		Amount          float64 `json:"-"`
		TransactionHash string  `json:"-"`
	}

	walletService interface {
		WithdrawRewards(ctx context.Context, userID uuid.UUID, amount float64) (string, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo rewardsRepository, ws walletService) *Service {
	return &Service{
		repo:      repo,
		ws:        ws,
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
	amount, err := s.repo.GetTotalAmount(ctx, uid)
	if err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("could not get total amount of rewards: %w", err)
	}

	txHash, err := s.ws.WithdrawRewards(ctx, uid, amount)
	if err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("could not create blockchain transaction: %w", err)
	}

	if err = s.repo.Withdraw(ctx, uid); err != nil {
		return ClaimRewardsResult{}, fmt.Errorf("ccould not update rewards status: %w", err)
	}

	return ClaimRewardsResult{
		Amount:          amount,
		DisplayAmount:   fmt.Sprintf("%.2f %s", amount, s.assetName),
		TransactionHash: txHash,
		TransactionURL:  fmt.Sprintf("https://explorer.solana.com/tx/%s?cluster=devnet", txHash),
	}, nil
}

// ClaimRewards send rewards to user by it and sets them to withdrawn.
func (s *Service) GetUserRewards(ctx context.Context, uid uuid.UUID) (float64, error) {
	amount, err := s.repo.GetTotalAmount(ctx, uid)
	if err != nil {
		return 0, fmt.Errorf("could not get total amount of rewards: %w", err)
	}
	return amount, nil
}
