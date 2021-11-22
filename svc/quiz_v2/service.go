package quiz_v2

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/consts"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/engine"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player/nats_player"
)

type (
	Service struct {
		engine *engine.Engine
	}
)

func NewService() *Service {
	s := &Service{
		engine: engine.New(),
	}

	return s
}

func (s *Service) GetQuizLink(ctx context.Context, uid uuid.UUID, username string, challengeID uuid.UUID) (*GetQuizLinkResponse, error) {
	baseQuizURL := nats.DefaultURL
	prefix := uid.String()
	recvMessageSubj := fmt.Sprintf("%v/%v", prefix, "recv")
	sendMessageSubj := fmt.Sprintf("%v/%v", prefix, "send")

	player, err := nats_player.NewNatsPlayer(uid.String(), consts.DefaultChallengeID, username, recvMessageSubj, sendMessageSubj)
	if err != nil {
		return nil, err
	}
	if err := player.Start(); err != nil {
		return nil, err
	}
	s.engine.AddPlayer(player)

	return &GetQuizLinkResponse{
		BaseQuizURL:     baseQuizURL,
		RecvMessageSubj: recvMessageSubj,
		SendMessageSubj: sendMessageSubj,
		UserID:          uid.String(),
	}, nil
}

func (s *Service) StartEngine() {
	s.engine.Start()
}
