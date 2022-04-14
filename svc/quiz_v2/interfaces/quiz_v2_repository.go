package interfaces

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"

	quiz_v2_repository "github.com/SatorNetwork/sator-api/svc/quiz_v2/repository"
)

type QuizV2Repository interface {
	RegisterNewPlayer(ctx context.Context, arg quiz_v2_repository.RegisterNewPlayerParams) error
	CountPlayersInRoom(ctx context.Context, challengeID uuid.UUID) (int64, error)
	UnregisterPlayer(ctx context.Context, arg quiz_v2_repository.UnregisterPlayerParams) error

	RegisterNewQuiz(ctx context.Context, arg quiz_v2_repository.RegisterNewQuizParams) (quiz_v2_repository.QuizzesV2, error)
	GetDistributedRewardsByChallengeID(ctx context.Context, challengeID uuid.UUID) (float64, error)
}
