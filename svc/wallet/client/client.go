package client

import (
	"context"

	"github.com/google/uuid"

	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	"github.com/SatorNetwork/sator-api/svc/wallet"
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
		GetListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID, limit, offset int32) (_ wallet.Transactions, err error)
		PayForService(ctx context.Context, uid uuid.UUID, amount float64, info string) error
		PayForNFT(ctx context.Context, uid uuid.UUID, amount float64, info string, creatorAddr string, creatorShare int32) error
		P2PTransfer(ctx context.Context, uid, recipientID uuid.UUID, amount float64, cfg *lib_solana.SendAssetsConfig, info string) error
		GetMultiplier(ctx context.Context, userID uuid.UUID) (_ int32, err error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// GetWalletsListByUserID ...
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

// GetListTransactionsByWalletID ...
func (c *Client) GetListTransactionsByWalletID(ctx context.Context, userID, walletID uuid.UUID, limit, offset int32) (_ wallet.Transactions, err error) {
	return c.s.GetListTransactionsByWalletID(ctx, userID, walletID, limit, offset)
}

// PayForService ...
func (c *Client) PayForService(ctx context.Context, uid uuid.UUID, amount float64, info string) error {
	return c.s.PayForService(ctx, uid, amount, info)
}

// PayForNFT ...
func (c *Client) PayForNFT(ctx context.Context, uid uuid.UUID, amount float64, info string, creatorAddr string, creatorShare int32) error {
	return c.s.PayForNFT(ctx, uid, amount, info, creatorAddr, creatorShare)
}

// P2PTransfer ...
func (c *Client) P2PTransfer(ctx context.Context, uid, recipientID uuid.UUID, amount float64, cfg *lib_solana.SendAssetsConfig, info string) error {
	return c.s.P2PTransfer(ctx, uid, recipientID, amount, cfg, info)
}

// GetMultiplier ...
func (c *Client) GetMultiplier(ctx context.Context, userID uuid.UUID) (_ int32, err error) {
	return c.s.GetMultiplier(ctx, userID)
}
