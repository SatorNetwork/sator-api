package gapi

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

type (
	SolanaClient struct {
		solana solanaClient
		wallet walletService
		mltpl  uint64

		tokenPubKey        string
		feeCollectorPubKey string
		feePayer           types.Account
		tokenPool          types.Account
	}

	solanaClient interface {
		GetTokenAccountBalanceWithAutoDerive(ctx context.Context, assetAddr, accountAddr string) (float64, error)
		CreateAccountWithATA(ctx context.Context, assetAddr, initAccAddr string, feePayer types.Account) (string, error)
		SendTransaction(ctx context.Context, feePayer, signer types.Account, instructions ...types.Instruction) (string, error)
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
		mltpl:              1e9,
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

	tx, err := c.sendAssetsWithAutoDerive(
		ctx,
		c.tokenPubKey,
		c.feePayer,
		c.tokenPool,
		walletAddr,
		amount,
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

	tx, err := c.sendAssetsWithAutoDerive(
		ctx,
		c.tokenPubKey,
		c.feePayer,
		userAcc,
		c.tokenPool.PublicKey.ToBase58(),
		amount,
	)
	if err != nil {
		return "", fmt.Errorf("could not claim rewards: %w", err)
	}

	return tx, nil
}

func (c *SolanaClient) sendAssetsWithAutoDerive(
	ctx context.Context,
	assetAddr string,
	feePayer, source types.Account,
	recipient string,
	amount float64) (string, error) {
	if _, err := c.solana.CreateAccountWithATA(ctx, assetAddr, recipient, feePayer); err != nil {
		log.Printf("CreateAccountWithATA: %v", err)
	}

	asset := common.PublicKeyFromString(assetAddr)

	sourceAta, _, err := common.FindAssociatedTokenAddress(
		common.PublicKeyFromString(source.PublicKey.ToBase58()),
		asset,
	)
	if err != nil {
		return "", fmt.Errorf("could not find associated token address: %w", err)
	}

	recipientAta, _, err := common.FindAssociatedTokenAddress(
		common.PublicKeyFromString(recipient),
		asset,
	)
	if err != nil {
		return "", fmt.Errorf("could not find associated token address: %w", err)
	}

	txHash, err := c.solana.SendTransaction(ctx, feePayer, source,
		tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
			From:     sourceAta,
			To:       recipientAta,
			Mint:     asset,
			Auth:     source.PublicKey,
			Signers:  []common.PublicKey{},
			Amount:   uint64(amount * float64(c.mltpl)),
			Decimals: 9,
		}),
	)
	if err != nil {
		return "", fmt.Errorf("could not send asset: %w", err)
	}

	return txHash, nil
}
