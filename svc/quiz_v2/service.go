package quiz_v2

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	internal_rsa "github.com/SatorNetwork/sator-api/internal/encryption/rsa"
	challenge_service "github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/SatorNetwork/sator-api/svc/profile"
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
		pc         profileClient

		serverRSAPrivateKey *rsa.PrivateKey
	}

	authClient interface {
		GetPublicKey(ctx context.Context, userID uuid.UUID) (*rsa.PublicKey, error)
	}

	profileClient interface {
		GetProfileByUserID(ctx context.Context, userID uuid.UUID, username string) (*profile.Profile, error)
	}
)

func NewService(
	natsURL,
	natsWSURL string,
	challenges interfaces.ChallengesService,
	stakeLevels interfaces.StakeLevels,
	rewards interfaces.RewardsService,
	ac authClient,
	pc profileClient,
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
		pc:                  pc,
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

	prefix := fmt.Sprintf("%v/%v", uuid.New().String(), uid.String())
	recvMessageSubj := fmt.Sprintf("%v/%v", prefix, "recv")
	sendMessageSubj := fmt.Sprintf("%v/%v", prefix, "send")

	publicKey, err := s.ac.GetPublicKey(ctx, uid)
	if err != nil {
		return nil, err
	}

	profile, err := s.pc.GetProfileByUserID(ctx, uid, username)
	if err != nil {
		return nil, errors.Wrap(err, "can't get profile by user id")
	}

	player, err := nats_player.NewNatsPlayer(
		uid.String(),
		challengeID.String(),
		username,
		profile.Avatar,
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

func (s *Service) GetChallengeByID(ctx context.Context, challengeID, userID uuid.UUID) (challenge_service.Challenge, error) {
	challenge, err := s.challenges.GetChallengeByID(ctx, challengeID, userID)
	if err != nil {
		return challenge_service.Challenge{}, errors.Wrap(err, "can't get challenge by ID")
	}

	roomDetails, err := s.engine.GetRoomDetails(challengeID.String())
	if err != nil {
		if _, ok := err.(*engine.ErrRoomNotFound); !ok {
			return challenge_service.Challenge{}, errors.Wrap(err, "can't get room details")
		} else {
			return challenge, nil
		}
	}

	challenge.Players = roomDetails.PlayersToStart
	challenge.RegisteredPlayers = roomDetails.RegisteredPlayers
	return challenge, nil
}
