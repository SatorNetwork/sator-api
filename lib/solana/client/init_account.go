//go:build !mock_solana

package client

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) InitAccountToUseAsset(ctx context.Context, feePayer, issuer, asset, initAcc types.Account) (string, error) {
	// Allow account hold the new asset
	rentExemptionBalanceIssuer, err := c.solana.GetMinimumBalanceForRentExemption(ctx, tokenprog.TokenAccountSize)
	if err != nil {
		return "", fmt.Errorf("could not get min balance for rent exemption: %w", err)
	}

	tx, err := c.SendTransaction(
		ctx,
		feePayer, initAcc,
		sysprog.CreateAccount(sysprog.CreateAccountParam{
			From:     feePayer.PublicKey,
			New:      initAcc.PublicKey,
			Owner:    common.TokenProgramID,
			Lamports: rentExemptionBalanceIssuer,
			Space:    tokenprog.TokenAccountSize,
		}),
		tokenprog.InitializeAccount(tokenprog.InitializeAccountParam{
			Account: initAcc.PublicKey,
			Mint:    asset.PublicKey,
			Owner:   issuer.PublicKey,
		}),
	)
	if err != nil {
		return "", fmt.Errorf("could not associate account: %w", err)
	}

	return tx, nil
}
