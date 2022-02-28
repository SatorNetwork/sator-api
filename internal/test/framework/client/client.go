package client

import (
	"log"

	"github.com/SatorNetwork/sator-api/internal/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/challenge"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/db"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/quiz_v2"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/shows"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/wallet"
)

type Client struct {
	Wallet           *wallet.WalletClient
	Auth             *auth.AuthClient
	QuizV2Client     *quiz_v2.QuizClient
	ChallengesClient *challenge.ChallengesClient
	ShowsClient      *shows.ShowsClient

	DB *db.DB
}

func NewClient() *Client {
	db, err := db.New()
	if err != nil {
		log.Fatalf("can't init DB: %v\n", err)
	}

	return &Client{
		Wallet:           wallet.New(),
		Auth:             auth.New(),
		QuizV2Client:     quiz_v2.New(),
		ChallengesClient: challenge.New(),
		ShowsClient:      shows.New(),

		DB: db,
	}
}
