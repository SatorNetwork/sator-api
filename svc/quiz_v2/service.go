package quiz_v2

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	internal_rsa "github.com/SatorNetwork/sator-api/lib/encryption/rsa"
	challenge_service "github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/SatorNetwork/sator-api/svc/profile"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/common"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/db/sql_builder"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/db/sql_executor"
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
		dbClient   *sql.DB
		qr         interfaces.QuizV2Repository

		serverRSAPrivateKey *rsa.PrivateKey

		disableRewardsForQuiz bool
		disableVerification   bool
	}

	authClient interface {
		GetPublicKey(ctx context.Context, userID uuid.UUID) (*rsa.PublicKey, error)
	}

	profileClient interface {
		GetProfileByUserID(ctx context.Context, userID uuid.UUID, username string) (*profile.Profile, error)
	}

	Challenge struct {
		ID               string `json:"id"`
		Title            string `json:"title"`
		PlayersToStart   int    `json:"players_to_start"`
		PlayersNumber    int    `json:"players_number"`
		PrizePool        string `json:"prize_pool"`
		IsRealmActivated bool   `json:"is_realm_activated"`
		Cover            string `json:"cover"`
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
	dbClient *sql.DB,
	qr interfaces.QuizV2Repository,
	serverRSAPrivateKey *rsa.PrivateKey,
	shuffleQuestions bool,
	disableRewardsForQuiz bool,
	quizLobbyLatency time.Duration,
	disableVerification bool,
) *Service {
	restrictionManager := restriction_manager.New(challenges)

	engine := engine.New(
		challenges,
		stakeLevels,
		rewards,
		qr,
		restrictionManager,
		shuffleQuestions,
		disableRewardsForQuiz,
		quizLobbyLatency,
	)

	s := &Service{
		engine:                engine,
		restrictionManager:    restrictionManager,
		natsURL:               natsURL,
		natsWSURL:             natsWSURL,
		challenges:            challenges,
		ac:                    ac,
		pc:                    pc,
		dbClient:              dbClient,
		qr:                    qr,
		serverRSAPrivateKey:   serverRSAPrivateKey,
		disableRewardsForQuiz: disableRewardsForQuiz,
		disableVerification:   disableVerification,
	}

	return s
}

func (s *Service) GetQuizLink(ctx context.Context, uid uuid.UUID, username string, challengeID uuid.UUID, mustAccess bool) (*GetQuizLinkResponse, error) {
	challenge, err := s.challenges.GetRawChallengeByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	questions, err := s.challenges.GetQuestionsByChallengeID(ctx, challengeID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get questions by challenge id")
	}
	if len(questions) < int(challenge.QuestionsPerGame) {
		return nil, errors.Errorf("can't choose %v questions out of %v", challenge.QuestionsPerGame, len(questions))
	}

	restricted, restrictionReason, err := s.restrictionManager.IsUserRestricted(ctx, challengeID, uid, mustAccess || s.disableVerification)
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

func (s *Service) GetChallengeByID(ctx context.Context, challengeID, userID uuid.UUID, mustAccess bool) (challenge_service.Challenge, error) {
	challenge, err := s.challenges.GetChallengeByID(ctx, challengeID, userID, mustAccess || s.disableVerification)
	if err != nil {
		return challenge_service.Challenge{}, errors.Wrap(err, "can't get challenge by ID")
	}

	currentPrizePool, err := common.GetCurrentPrizePool(
		s.qr,
		challenge.ID,
		challenge.PrizePoolAmount,
		challenge.MinimumReward,
		challenge.PercentForQuiz,
	)
	if err != nil {
		return challenge_service.Challenge{}, err
	}

	if !s.disableRewardsForQuiz {
		challenge.CurrentPrizePool = fmt.Sprintf("%v SAO", currentPrizePool)
	}
	challenge.CurrentPrizePoolAmount = currentPrizePool

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
	challenge.RegisteredPlayersInDB = roomDetails.RegisteredPlayersInDB
	return challenge, nil
}

func (s *Service) GetChallengesSortedByPlayers(ctx context.Context, userID uuid.UUID, limit, offset int32, mustAccess bool) ([]*Challenge, error) {
	query := sql_builder.ConstructGetChallengesSortedByPlayersQuery(userID, limit, offset)
	{
		tmpl := `
		Challenges sorted by players params:
		User's ID: %v
		Limit:     %v
		Offset:    %v
		`
		log.Printf(tmpl, userID, limit, offset)
		log.Printf("Challenges sorted by players query: %v", query)
	}
	sqlExecutor := sql_executor.New(s.dbClient)
	sqlChallenges, err := sqlExecutor.ExecuteGetChallengesSortedByPlayersQuery(query, nil)
	if err != nil {
		return nil, err
	}
	{
		serializedChallenges, err := json.Marshal(sqlChallenges)
		if err != nil {
			return nil, err
		}
		log.Printf("Serialized SQL challenges: %s\n", serializedChallenges)
	}
	challenges, err := s.NewChallengesFromSQL(ctx, sqlChallenges, userID, mustAccess || s.disableVerification)
	if err != nil {
		return nil, err
	}
	{
		serializedChallenges, err := json.Marshal(challenges)
		if err != nil {
			return nil, err
		}
		log.Printf("Serialized challenges: %s\n", serializedChallenges)
	}

	return challenges, nil
}

func (s *Service) NewChallengeFromSQL(ctx context.Context, c *sql_executor.Challenge, userID uuid.UUID, mustAccess bool) (*Challenge, error) {
	var currentPrizePool float64
	{
		challengeUID, err := uuid.Parse(c.ID)
		if err != nil {
			return nil, err
		}

		// TODO(evg): get rid of extra database call
		challenge, err := s.challenges.GetChallengeByID(ctx, challengeUID, userID, mustAccess || s.disableVerification)
		if err != nil {
			return nil, errors.Wrap(err, "can't get challenge by ID")
		}

		currentPrizePool, err = common.GetCurrentPrizePool(
			s.qr,
			challenge.ID,
			challenge.PrizePoolAmount,
			challenge.MinimumReward,
			challenge.PercentForQuiz,
		)
		if err != nil {
			return nil, err
		}
	}

	isActivated := c.IsActivated
	if mustAccess || s.disableVerification {
		isActivated = true
	}

	return &Challenge{
		ID:               c.ID,
		Title:            c.Title,
		PlayersToStart:   c.PlayersToStart,
		PlayersNumber:    c.PlayersNum,
		PrizePool:        fmt.Sprintf("%.2f SAO", currentPrizePool),
		IsRealmActivated: isActivated,
		Cover:            c.Cover.String,
	}, nil
}

func (s *Service) NewChallengesFromSQL(ctx context.Context, sqlChallenges []*sql_executor.Challenge, userID uuid.UUID, mustAccess bool) ([]*Challenge, error) {
	challenges := make([]*Challenge, 0)
	for _, c := range sqlChallenges {
		challenge, err := s.NewChallengeFromSQL(ctx, c, userID, mustAccess || s.disableVerification)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, challenge)
	}

	return challenges, nil
}
