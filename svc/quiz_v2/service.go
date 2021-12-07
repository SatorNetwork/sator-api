package quiz_v2

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/internal/db"
	quiz_v2_challenge "github.com/SatorNetwork/sator-api/svc/quiz_v2/challenge"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/engine"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player/nats_player"
)

type (
	Service struct {
		engine *engine.Engine

		natsURL    string
		challenges quiz_v2_challenge.ChallengesService
	}
)

func NewService(natsURL string, challenges quiz_v2_challenge.ChallengesService) *Service {
	s := &Service{
		engine:     engine.New(challenges),
		natsURL:    natsURL,
		challenges: challenges,
	}

	return s
}

func (s *Service) GetQuizLink(ctx context.Context, uid uuid.UUID, username string, challengeID uuid.UUID) (*GetQuizLinkResponse, error) {
	receivedReward, err := s.challenges.GetChallengeReceivedRewardAmountByUserID(ctx, challengeID, uid)
	if err != nil && !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("could not get received reward amount: %w", err)
	}
	if receivedReward > 0 {
		return nil, errors.New("reward has been already received for this challenge")
	}

	challengeByID, err := s.challenges.GetChallengeByID(ctx, challengeID, uid)
	if err != nil {
		return nil, fmt.Errorf("could not found challenge: %w", err)
	}
	attempts, err := s.challenges.GetPassedChallengeAttempts(ctx, challengeID, uid)
	if err != nil {
		return nil, fmt.Errorf("could not get passed challenge attempts: %w", err)
	}
	if attempts >= int64(challengeByID.UserMaxAttempts) {
		return nil, errors.New("no more attempts left")
	}

	prefix := uid.String()
	recvMessageSubj := fmt.Sprintf("%v/%v", prefix, "recv")
	sendMessageSubj := fmt.Sprintf("%v/%v", prefix, "send")

	player, err := nats_player.NewNatsPlayer(uid.String(), challengeID.String(), username, s.natsURL, recvMessageSubj, sendMessageSubj)
	if err != nil {
		return nil, errors.Wrap(err, "can't create nats player")
	}
	if err := player.Start(); err != nil {
		return nil, err
	}
	s.engine.AddPlayer(player)

	return &GetQuizLinkResponse{
		BaseQuizURL:     s.natsURL,
		RecvMessageSubj: recvMessageSubj,
		SendMessageSubj: sendMessageSubj,
		UserID:          uid.String(),
	}, nil
}

func (s *Service) StartEngine() {
	s.engine.Start()
}
