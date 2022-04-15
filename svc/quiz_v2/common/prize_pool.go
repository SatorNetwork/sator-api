package common

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
)

func GetCurrentPrizePool(
	qr interfaces.QuizV2Repository,
	challengeID uuid.UUID,
	prizePoolAmount float64,
	minimumReward float64,
	percentForQuiz float64,
) (float64, error) {
	ctxb := context.Background()
	distributedRewards, err := qr.GetDistributedRewardsByChallengeID(ctxb, challengeID)
	if err != nil && !strings.Contains(err.Error(), "converting NULL to float64 is unsupported") {
		return 0, errors.Wrap(err, "can't get distributed rewards by challenge id")
	}
	if err != nil && strings.Contains(err.Error(), "converting NULL to float64 is unsupported") {
		distributedRewards = 0
	}
	leftInPool := prizePoolAmount - distributedRewards
	if leftInPool <= 0 {
		return 0, errors.Wrap(err, "no money left in pool")
	}
	if leftInPool <= minimumReward {
		return leftInPool, nil
	}

	currentPrizePool := leftInPool / 100 * percentForQuiz
	if currentPrizePool < minimumReward {
		currentPrizePool = minimumReward
	}

	return currentPrizePool, nil
}
