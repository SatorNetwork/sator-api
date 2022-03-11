//go:build !mock_solana

package client

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/tokenprog"
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
		sysprog.CreateAccount(
			feePayer.PublicKey,
			initAcc.PublicKey,
			common.TokenProgramID,
			rentExemptionBalanceIssuer,
			tokenprog.TokenAccountSize,
		),
		tokenprog.InitializeAccount(
			initAcc.PublicKey,
			asset.PublicKey,
			issuer.PublicKey,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not associate account: %w", err)
	}

	return tx, nil
}
