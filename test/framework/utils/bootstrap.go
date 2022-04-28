package utils

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"
	"github.com/stretchr/testify/require"

	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
	exchange_rates_client "github.com/SatorNetwork/sator-api/svc/exchange_rates/client"
	"github.com/SatorNetwork/sator-api/test/framework/accounts"
	"github.com/SatorNetwork/sator-api/test/framework/client"
)

var bootstrapLock = &sync.Mutex{}

func BootstrapIfNeeded(ctx context.Context, t *testing.T) error {
	bootstrapLock.Lock()
	defer bootstrapLock.Unlock()

	c := client.NewClient()

	exchangeRatesClient, err := exchange_rates_client.Easy(c.DB.Client())
	require.NoError(t, err)
	sc := solana_client.New("http://localhost:8899", solana_client.Config{
		SystemProgram:  common.SystemProgramID.ToBase58(),
		SysvarRent:     common.SysVarRentPubkey.ToBase58(),
		SysvarClock:    common.SysVarClockPubkey.ToBase58(),
		SplToken:       common.TokenProgramID.ToBase58(),
		StakeProgramID: "CL9tjeJL38C3eWqd6g7iHMnXaJ17tmL2ygkLEHghrj4u",
	}, exchangeRatesClient)

	if err := c.DB.PuzzleGameDB().Bootstrap(ctx); err != nil {
		return err
	}

	needed, err := CheckIfBootstrapNeeded(ctx, sc)
	if err != nil {
		return errors.Wrap(err, "can't check if bootstrap is needed")
	}
	if !needed {
		return nil
	}

	if err := Bootstrap(ctx, t, sc); err != nil {
		return err
	}

	if err := c.DB.Bootstrap(ctx); err != nil {
		return err
	}

	return nil
}

func CheckIfBootstrapNeeded(ctx context.Context, sc lib_solana.Interface) (bool, error) {
	_, tokenHolder, asset := accounts.GetAccounts()

	tokenHolderAta, _, err := common.FindAssociatedTokenAddress(tokenHolder.PublicKey, asset.PublicKey)
	if err != nil {
		return false, err
	}
	balance, err := sc.GetTokenAccountBalance(ctx, tokenHolderAta.ToBase58())
	if err != nil && strings.Contains(err.Error(), `{"code":-32602,"message":"Invalid param: could not find account"}`) {
		return true, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "can't get token account balance for token holder")
	}

	return balance == 0, nil
}

func Bootstrap(ctx context.Context, t *testing.T, sc lib_solana.Interface) error {
	airdropSolToFeePayer(ctx, t, sc)
	createAsset(ctx, t, sc)
	issueTokensToTokenHolder(ctx, t, sc)

	return nil
}

func airdropSolToFeePayer(ctx context.Context, t *testing.T, solanaClient lib_solana.Interface) {
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

func createAsset(ctx context.Context, t *testing.T, solanaClient lib_solana.Interface) {
	feePayer, tokenHolder, asset := accounts.GetAccounts()

	_, err := solanaClient.CreateAsset(
		ctx,
		feePayer,
		tokenHolder,
		asset,
	)
	require.NoError(t, err)
}

func issueTokensToTokenHolder(ctx context.Context, t *testing.T, solanaClient lib_solana.Interface) {
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
			feePayer,
			tokenHolder,
			asset,
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
