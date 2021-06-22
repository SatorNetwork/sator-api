package client

import (
	"context"

	"github.com/SatorNetwork/sator-api/svc/wallet"

	"github.com/google/uuid"
)

type (
	// Client struct
	Client struct {
		s service
	}

	service interface {
		GetWallets(ctx context.Context, userID uuid.UUID) (wallet.Wallets, error)
		GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (wallet.Wallet, error)
		CreateWallet(ctx context.Context, userID uuid.UUID) error
		WithdrawRewards(ctx context.Context, userID uuid.UUID, amount float64) (tx string, err error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// GetWalletsList ...
func (c *Client) GetWalletsListByUserID(ctx context.Context, userID uuid.UUID) (wallet.Wallets, error) {
	return c.s.GetWallets(ctx, userID)
}

// GetWalletByID ...
func (c *Client) GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (wallet.Wallet, error) {
	return c.s.GetWalletByID(ctx, userID, walletID)
}

// CreateWallet ...
func (c *Client) CreateWallet(ctx context.Context, userID uuid.UUID) error {
	return c.s.CreateWallet(ctx, userID)
}

// WithdrawRewards ...
func (c *Client) WithdrawRewards(ctx context.Context, userID uuid.UUID, amount float64) (tx string, err error) {
	return c.s.WithdrawRewards(ctx, userID, amount)
}
