package restriction_manager

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/internal/db"
	challenge_service "github.com/SatorNetwork/sator-api/svc/challenge"
)

type restrictionReason uint8

const (
	undefined restrictionReason = iota
	rewardAlreadyReceived
	noMoreAttemptsLeft
)

func (r restrictionReason) String() string {
	switch r {
	case undefined:
		return "undefined"
	case rewardAlreadyReceived:
		return "reward has been already received for this challenge"
	case noMoreAttemptsLeft:
		return "no more attempts left"
	default:
		return "undefined"
	}
}

type (
	RestrictionManager interface {
		IsUserRestricted(ctx context.Context, challengeID, userID uuid.UUID) (bool, restrictionReason, error)
		RegisterEarnedReward(ctx context.Context, challengeID, userID uuid.UUID, rewardAmount float64) error
	}

	restrictionManager struct {
		challenge challengeService
	}

	challengeService interface {
		GetChallengeByID(ctx context.Context, challengeID, userID uuid.UUID) (challenge_service.Challenge, error)
		GetChallengeReceivedRewardAmountByUserID(ctx context.Context, challengeID, userID uuid.UUID) (float64, error)
		GetPassedChallengeAttempts(ctx context.Context, challengeID, userID uuid.UUID) (int64, error)
		StoreChallengeReceivedRewardAmount(ctx context.Context, challengeID, userID uuid.UUID, rewardAmount float64) error
	}
)

func New(challenge challengeService) RestrictionManager {
	return &restrictionManager{
		challenge: challenge,
	}
}

func (m *restrictionManager) IsUserRestricted(ctx context.Context, challengeID, userID uuid.UUID) (bool, restrictionReason, error) {
	receivedReward, err := m.challenge.GetChallengeReceivedRewardAmountByUserID(ctx, challengeID, userID)
	if err != nil && !db.IsNotFoundError(err) {
		return false, undefined, errors.Wrap(err, "could not get received reward amount")
	}
	if receivedReward > 0 {
		return true, rewardAlreadyReceived, nil
	}

	challenge, err := m.challenge.GetChallengeByID(ctx, challengeID, userID)
	if err != nil {
		return false, undefined, errors.Wrap(err, "can't get challenge by ID")
	}
	attempts, err := m.challenge.GetPassedChallengeAttempts(ctx, challengeID, userID)
	if err != nil {
		return false, undefined, errors.Wrap(err, "could not get passed challenge attempts")
	}
	if attempts >= int64(challenge.UserMaxAttempts) {
		return true, noMoreAttemptsLeft, nil
	}

	return false, undefined, nil
}

func (m *restrictionManager) RegisterEarnedReward(ctx context.Context, challengeID, userID uuid.UUID, rewardAmount float64) error {
	return m.challenge.StoreChallengeReceivedRewardAmount(ctx, challengeID, userID, rewardAmount)
}
