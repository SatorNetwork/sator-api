package interfaces

import (
	"context"

	"github.com/google/uuid"

	"github.com/SatorNetwork/sator-api/svc/challenge"
)

type ChallengesService interface {
	GetChallengeByID(ctx context.Context, challengeID, userID uuid.UUID) (challenge.Challenge, error)
	GetRawChallengeByID(ctx context.Context, challengeID uuid.UUID) (challenge.RawChallenge, error)
	GetQuestionsByChallengeID(ctx context.Context, challengeID uuid.UUID) ([]challenge.Question, error)
	GetChallengeReceivedRewardAmountByUserID(ctx context.Context, challengeID, userID uuid.UUID) (float64, error)
	GetPassedChallengeAttempts(ctx context.Context, challengeID, userID uuid.UUID) (int64, error)
	StoreChallengeReceivedRewardAmount(ctx context.Context, challengeID, userID uuid.UUID, rewardAmount float64) error
}
