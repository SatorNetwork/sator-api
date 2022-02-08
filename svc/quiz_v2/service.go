package quiz_v2

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/internal/db"
	internal_rsa "github.com/SatorNetwork/sator-api/internal/encryption/rsa"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/engine"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player/nats_player"
)

type (
	Service struct {
		engine *engine.Engine

		natsURL    string
		natsWSURL  string
		challenges interfaces.ChallengesService
		ac         authClient

		serverRSAPrivateKey *rsa.PrivateKey
	}

	authClient interface {
		GetPublicKey(ctx context.Context, userID uuid.UUID) (*rsa.PublicKey, error)
	}
)

func NewService(
	natsURL,
	natsWSURL string,
	challenges interfaces.ChallengesService,
	stakeLevels interfaces.StakeLevels,
	rewards interfaces.RewardsService,
	ac authClient,
	serverRSAPrivateKey *rsa.PrivateKey,
	shuffleQuestions bool,
) *Service {
	s := &Service{
		engine:              engine.New(challenges, stakeLevels, rewards, shuffleQuestions),
		natsURL:             natsURL,
		natsWSURL:           natsWSURL,
		challenges:          challenges,
		ac:                  ac,
		serverRSAPrivateKey: serverRSAPrivateKey,
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

	publicKey, err := s.ac.GetPublicKey(ctx, uid)
	if err != nil {
		return nil, err
	}

	player, err := nats_player.NewNatsPlayer(
		uid.String(),
		challengeID.String(),
		username,
		s.natsURL,
		recvMessageSubj,
		sendMessageSubj,
		publicKey,
		s.serverRSAPrivateKey,
	)
	if err != nil {
		return nil, errors.Wrap(err, "can't create nats player")
	}
	if err := player.Start(); err != nil {
		return nil, err
	}
	s.engine.AddPlayer(player)

	publicKeyBytes, err := internal_rsa.PublicKeyToBytes(&s.serverRSAPrivateKey.PublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "can't encode server's public key")
	}

	return &GetQuizLinkResponse{
		BaseQuizURL:     s.natsURL,
		BaseQuizWSURL:   s.natsWSURL,
		RecvMessageSubj: recvMessageSubj,
		SendMessageSubj: sendMessageSubj,
		UserID:          uid.String(),
		ServerPublicKey: string(publicKeyBytes),
	}, nil
}

func (s *Service) StartEngine() {
	s.engine.Start()
}
