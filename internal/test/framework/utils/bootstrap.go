package utils

import (
	"context"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/internal/solana"
	"github.com/SatorNetwork/sator-api/internal/test/framework/accounts"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client"
)

var bootstrapLock = &sync.Mutex{}

func BootstrapIfNeeded(ctx context.Context, t *testing.T) error {
	bootstrapLock.Lock()
	defer bootstrapLock.Unlock()

	needed, err := CheckIfBootstrapNeeded(ctx)
	if err != nil {
		return err
	}
	if !needed {
		return nil
	}

	if err := Bootstrap(ctx, t); err != nil {
		return err
	}

	c := client.NewClient()
	if err := c.DB.Bootstrap(ctx); err != nil {
		return err
	}

	return nil
}

func CheckIfBootstrapNeeded(ctx context.Context) (bool, error) {
	sc := solana.New("http://localhost:8899", solana.Config{
		SystemProgram:  "11111111111111111111111111111111",
		SysvarRent:     "SysvarRent111111111111111111111111111111111",
		SysvarClock:    "SysvarC1ock11111111111111111111111111111111",
		SplToken:       "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		StakeProgramID: "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u",
	})
	_, tokenHolder, asset := accounts.GetAccounts()

	tokenHolderAta, _, err := common.FindAssociatedTokenAddress(tokenHolder.PublicKey, asset.PublicKey)
	if err != nil {
		return false, err
	}
	balance, err := sc.GetTokenAccountBalance(ctx, tokenHolderAta.ToBase58())
	if err != nil {
		return false, err
	}

	return balance == 0, nil
}

func Bootstrap(ctx context.Context, t *testing.T) error {
	airdropSolToFeePayer(ctx, t)
	createAsset(ctx, t)
	issueTokensToTokenHolder(ctx, t)

	return nil
}

func airdropSolToFeePayer(ctx context.Context, t *testing.T) {
	solanaClient := solana.New("http://localhost:8899", solana.Config{
		SystemProgram:  "11111111111111111111111111111111",
		SysvarRent:     "SysvarRent111111111111111111111111111111111",
		SysvarClock:    "SysvarC1ock11111111111111111111111111111111",
		SplToken:       "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		StakeProgramID: "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u",
	})
	feePayer := accounts.GetFeePayer()
	const solToAirdrop = 1

	BackoffRetry(t, func() error {
		_, err := solanaClient.RequestAirdrop(ctx, feePayer.PublicKey.ToBase58(), solToAirdrop)
		return err
	})

	BackoffRetry(t, func() error {
		balance, err := solanaClient.GetAccountBalanceSOL(ctx, feePayer.PublicKey.ToBase58())
		require.NoError(t, err)
		if balance != solToAirdrop {
			return errors.Errorf("unexpected account balance SOL, want: %v, got: %v", solToAirdrop, balance)
		}

		return nil
	})
}

func createAsset(ctx context.Context, t *testing.T) {
	solanaClient := solana.New("http://localhost:8899", solana.Config{
		SystemProgram:  "11111111111111111111111111111111",
		SysvarRent:     "SysvarRent111111111111111111111111111111111",
		SysvarClock:    "SysvarC1ock11111111111111111111111111111111",
		SplToken:       "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		StakeProgramID: "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u",
	})
	feePayer, tokenHolder, asset := accounts.GetAccounts()

	_, err := solanaClient.CreateAsset(
		ctx,
		solanaClient.AccountFromPrivateKeyBytes(feePayer.PrivateKey),
		solanaClient.AccountFromPrivateKeyBytes(tokenHolder.PrivateKey),
		solanaClient.AccountFromPrivateKeyBytes(asset.PrivateKey),
	)
	require.NoError(t, err)
}

func issueTokensToTokenHolder(ctx context.Context, t *testing.T) {
	solanaClient := solana.New("http://localhost:8899", solana.Config{
		SystemProgram:  "11111111111111111111111111111111",
		SysvarRent:     "SysvarRent111111111111111111111111111111111",
		SysvarClock:    "SysvarC1ock11111111111111111111111111111111",
		SplToken:       "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
		StakeProgramID: "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u",
	})
	feePayer, tokenHolder, asset := accounts.GetAccounts()
	const tokensToIssue = 500000000

	tokenHolderAta, _, err := common.FindAssociatedTokenAddress(tokenHolder.PublicKey, asset.PublicKey)
	require.NoError(t, err)

	BackoffRetry(t, func() error {
		_, err := solanaClient.CreateAccountWithATA(ctx, asset.PublicKey.ToBase58(), tokenHolder.PublicKey.ToBase58(), feePayer)
		return err
	})

	BackoffRetry(t, func() error {
		_, err := solanaClient.IssueAsset(
			ctx,
			solanaClient.AccountFromPrivateKeyBytes(feePayer.PrivateKey),
			solanaClient.AccountFromPrivateKeyBytes(tokenHolder.PrivateKey),
			solanaClient.AccountFromPrivateKeyBytes(asset.PrivateKey),
			tokenHolderAta,
			tokensToIssue,
		)
		return err
	})

	BackoffRetry(t, func() error {
		balance, err := solanaClient.GetTokenAccountBalance(context.Background(), tokenHolderAta.ToBase58())
		require.NoError(t, err)

		if balance != tokensToIssue {
			return errors.Errorf("unexpected token account balance, want: %v, got: %v", tokensToIssue, balance)
		}

		return nil
	})
}
