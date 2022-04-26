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

func (c *Client) CreateAsset(ctx context.Context, feePayer, issuer, asset types.Account) (string, error) {
	rentExemptionBalance, err := c.solana.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		return "", fmt.Errorf("could not get min balance for rent exemption: %w", err)
	}
	// Transform general account to asset
	tx, err := c.SendTransaction(
		ctx,
		feePayer, asset,
		sysprog.CreateAccount(sysprog.CreateAccountParam{
			From:     feePayer.PublicKey,
			New:      asset.PublicKey,
			Owner:    common.TokenProgramID,
			Lamports: rentExemptionBalance,
			Space:    tokenprog.MintAccountSize,
		}),
		tokenprog.InitializeMint(tokenprog.InitializeMintParam{
			Decimals:   c.decimals,
			Mint:       asset.PublicKey,
			MintAuth:   issuer.PublicKey,
			FreezeAuth: nil,
		}),
	)
	if err != nil {
		return "", fmt.Errorf("could not issue new asset: %w", err)
	}

	return tx, nil
}
