package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) IssueAsset(ctx context.Context, dest types.Account, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))
	// Issue Asset
	tx, err := c.SendTransaction(
		ctx, c.Issuer,
		tokenprog.MintToChecked(
			c.Asset.PublicKey,
			dest.PublicKey,
			c.Issuer.PublicKey,
			[]common.PublicKey{},
			amountToSend,
			c.decimals,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not issue additional Asset amount: %w", err)
	}
	return tx, nil
}
