package quiz_v2

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	internal_rsa "github.com/SatorNetwork/sator-api/internal/encryption/rsa"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/engine"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player/nats_player"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/restriction_manager"
)

type (
	Service struct {
		engine             *engine.Engine
		restrictionManager restriction_manager.RestrictionManager

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
	restrictionManager := restriction_manager.New(challenges)

	s := &Service{
		engine:              engine.New(challenges, stakeLevels, rewards, restrictionManager, shuffleQuestions),
		restrictionManager:  restrictionManager,
		natsURL:             natsURL,
		natsWSURL:           natsWSURL,
		challenges:          challenges,
		ac:                  ac,
		serverRSAPrivateKey: serverRSAPrivateKey,
	}

	return s
}

func (s *Service) GetQuizLink(ctx context.Context, uid uuid.UUID, username string, challengeID uuid.UUID) (*GetQuizLinkResponse, error) {
	restricted, restrictionReason, err := s.restrictionManager.IsUserRestricted(ctx, challengeID, uid)
	if err != nil {
		return nil, errors.Wrap(err, "can't check if user is restricted")
	}
	if restricted {
		return nil, errors.Errorf("user is restricted for this challenge reason: %v", restrictionReason.String())
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
