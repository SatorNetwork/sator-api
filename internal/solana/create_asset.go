package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/tokenprog"
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
		sysprog.CreateAccount(
			feePayer.PublicKey,
			asset.PublicKey,
			common.TokenProgramID,
			rentExemptionBalance,
			tokenprog.MintAccountSize,
		),
		tokenprog.InitializeMint(
			c.decimals,
			asset.PublicKey,
			issuer.PublicKey,
			common.PublicKey{},
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not issue new asset: %w", err)
	}

	return tx, nil
}
