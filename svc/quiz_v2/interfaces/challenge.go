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
}

type StaticChallenges struct{}

func (mock *StaticChallenges) GetChallengeByID(ctx context.Context, challengeID, userID uuid.UUID) (challenge.Challenge, error) {
	return challenge.Challenge{}, nil
}

func (mock *StaticChallenges) GetRawChallengeByID(ctx context.Context, challengeID uuid.UUID) (challenge.RawChallenge, error) {
	return challenge.RawChallenge{}, nil
}

func (mock *StaticChallenges) GetQuestionsByChallengeID(ctx context.Context, challengeID uuid.UUID) ([]challenge.Question, error) {
	return []challenge.Question{}, nil
}

func (mock *StaticChallenges) GetChallengeReceivedRewardAmountByUserID(ctx context.Context, challengeID, userID uuid.UUID) (float64, error) {
	return 1, nil
}

func (mock *StaticChallenges) GetPassedChallengeAttempts(ctx context.Context, challengeID, userID uuid.UUID) (int64, error) {
	return 1, nil
}
