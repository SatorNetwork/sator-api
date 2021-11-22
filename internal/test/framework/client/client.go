package client

import (
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/quiz_v2"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/wallet"
)

type Client struct {
	Wallet       *wallet.WalletClient
	Auth         *auth.AuthClient
	QuizV2Client *quiz_v2.QuizClient
}

func NewClient() *Client {
	return &Client{
		Wallet:       wallet.NewWalletClient(),
		Auth:         auth.NewAuthClient(),
		QuizV2Client: quiz_v2.NewQuizClient(),
	}
}
