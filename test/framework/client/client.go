package client

import (
	"log"

	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/client/challenge"
	"github.com/SatorNetwork/sator-api/test/framework/client/db"
	"github.com/SatorNetwork/sator-api/test/framework/client/puzzle_game"
	"github.com/SatorNetwork/sator-api/test/framework/client/quiz_v2"
	"github.com/SatorNetwork/sator-api/test/framework/client/shows"
	"github.com/SatorNetwork/sator-api/test/framework/client/trading_platforms"
	"github.com/SatorNetwork/sator-api/test/framework/client/wallet"
)

type Client struct {
	Wallet                 *wallet.WalletClient
	Auth                   *auth.AuthClient
	QuizV2Client           *quiz_v2.QuizClient
	ChallengesClient       *challenge.ChallengesClient
	ShowsClient            *shows.ShowsClient
	TradingPlatformsClient *trading_platforms.TradingPlatformsClient
	PuzzleGameClient       *puzzle_game.PuzzleGameClient

	DB *db.DB
}

func NewClient() *Client {
	db, err := db.New()
	if err != nil {
		log.Fatalf("can't init DB: %v\n", err)
	}

	return &Client{
		Wallet:                 wallet.New(),
		Auth:                   auth.New(),
		QuizV2Client:           quiz_v2.New(),
		ChallengesClient:       challenge.New(),
		ShowsClient:            shows.New(),
		TradingPlatformsClient: trading_platforms.New(),

		DB: db,
	}
}
