package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) IssueAsset(ctx context.Context, feePayer, issuer, asset types.Account, dest common.PublicKey, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))
	// Issue asset
	tx, err := c.SendTransaction(
		ctx,
		feePayer, issuer,
		tokenprog.MintToChecked(
			asset.PublicKey,
			dest,
			issuer.PublicKey,
			[]common.PublicKey{feePayer.PublicKey, issuer.PublicKey},
			amountToSend,
			c.decimals,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not issue additional asset amount: %w", err)
	}
	return tx, nil
}
