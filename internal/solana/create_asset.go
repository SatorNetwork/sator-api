package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/tokenprog"
)

func (c *Client) CreateAsset(ctx context.Context) (string, error) {
	rentExemptionBalance, err := c.solana.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		return "", fmt.Errorf("could not get min balance for rent exemption: %w", err)
	}
	// Transform general account to Asset
	tx, err := c.SendTransaction(
		ctx, c.Asset,
		sysprog.CreateAccount(
			c.FeePayer.PublicKey,
			c.Asset.PublicKey,
			common.TokenProgramID,
			rentExemptionBalance,
			tokenprog.MintAccountSize,
		),
		tokenprog.InitializeMint(
			c.decimals,
			c.Asset.PublicKey,
			c.Issuer.PublicKey,
			common.PublicKey{},
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not issue new Asset: %w", err)
	}

	return tx, nil
}
