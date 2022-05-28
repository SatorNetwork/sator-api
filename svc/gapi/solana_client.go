package gapi

import (
	"context"
	"fmt"
	"log"

	"github.com/SatorNetwork/sator-api/lib/solana"
	"github.com/google/uuid"
	"github.com/portto/solana-go-sdk/types"
)

type (
	SolanaClient struct {
		solana solanaClient
		wallet walletService

		tokenPubKey        string
		feeCollectorPubKey string
		feePayer           types.Account
		tokenPool          types.Account
	}

	solanaClient interface {
		GetTokenAccountBalanceWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) (float64, error)
		SendAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, source types.Account, recipientAddr string, amount float64, cfg *solana.SendAssetsConfig) (string, error)
	}

	walletService interface {
		GetUserSolanaAccount(ctx context.Context, userID uuid.UUID) ([]byte, error)
	}
)

// NewSolanaClient ...
func NewSolanaClient(solana solanaClient, wallet walletService, tokenPubKey, feeCollectorPubKey string, feePayer, tokenPool types.Account) *SolanaClient {
	return &SolanaClient{
		solana:             solana,
		wallet:             wallet,
		tokenPubKey:        tokenPubKey,
		feeCollectorPubKey: feeCollectorPubKey,
		feePayer:           feePayer,
		tokenPool:          tokenPool,
	}
}

func (c *SolanaClient) GetUserWalletAddress(ctx context.Context, uid uuid.UUID) (string, error) {
	userSolAcc, err := c.wallet.GetUserSolanaAccount(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("get user solana account: %w", err)
	}

	solAcc, err := types.AccountFromBytes(userSolAcc)
	if err != nil {
		return "", fmt.Errorf("parse user solana account: %w", err)
	}

	return solAcc.PublicKey.ToBase58(), nil
}

func (c *SolanaClient) GetUserSolanaAccount(ctx context.Context, uid uuid.UUID) (types.Account, error) {
	userSolAcc, err := c.wallet.GetUserSolanaAccount(ctx, uid)
	if err != nil {
		return types.Account{}, fmt.Errorf("get user solana account: %w", err)
	}

	solAcc, err := types.AccountFromBytes(userSolAcc)
	if err != nil {
		return types.Account{}, fmt.Errorf("parse user solana account: %w", err)
	}

	return solAcc, nil
}

func (c *SolanaClient) GetBalance(ctx context.Context, uid uuid.UUID) (float64, error) {
	log.Println("get balance", uid)

	walletAddr, err := c.GetUserWalletAddress(ctx, uid)
	if err != nil {
		return 0, fmt.Errorf("get user wallet address: %w", err)
	}

	balance, err := c.solana.GetTokenAccountBalanceWithAutoDerive(ctx, c.tokenPubKey, walletAddr)
	if err != nil {
		return 0, fmt.Errorf("get token account balance: %w", err)
	}

	return balance, nil
}

func (c *SolanaClient) ClaimRewards(ctx context.Context, uid uuid.UUID, amount float64) (string, error) {
	log.Println("claim rewards", uid, amount)

	walletAddr, err := c.GetUserWalletAddress(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("get user wallet address: %w", err)
	}

	tx, err := c.solana.SendAssetsWithAutoDerive(
		ctx,
		c.tokenPubKey,
		c.feePayer,
		c.tokenPool,
		walletAddr,
		amount,
		&solana.SendAssetsConfig{},
	)
	if err != nil {
		return "", fmt.Errorf("could not claim rewards: %w", err)
	}

	return tx, nil
}

func (c *SolanaClient) Pay(ctx context.Context, uid uuid.UUID, amount float64, info string) (string, error) {
	log.Println("pay", uid, amount, info)

	userAcc, err := c.GetUserSolanaAccount(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("get user solana account: %w", err)
	}

	tx, err := c.solana.SendAssetsWithAutoDerive(
		ctx,
		c.tokenPubKey,
		c.feePayer,
		userAcc,
		c.tokenPool.PublicKey.ToBase58(),
		amount,
		&solana.SendAssetsConfig{},
	)
	if err != nil {
		return "", fmt.Errorf("could not claim rewards: %w", err)
	}

	return tx, nil
}
