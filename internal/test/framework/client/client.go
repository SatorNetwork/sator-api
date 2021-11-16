package client

import (
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/wallet"
)

type Client struct {
	Wallet *wallet.WalletClient
	Auth   *auth.AuthClient
}

func NewClient() *Client {
	return &Client{
		Wallet: wallet.NewWalletClient(),
		Auth:   auth.NewAuthClient(),
	}
}
