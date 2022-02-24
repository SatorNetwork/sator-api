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
}
