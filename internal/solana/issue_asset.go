package solana

import (
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) IssueAsset(feePayer, issuer, asset, dest types.Account, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))
	// Issue asset
	tx, err := c.SendTransaction(
		feePayer, issuer,
		tokenprog.MintToChecked(
			asset.PublicKey,
			dest.PublicKey,
			issuer.PublicKey,
			[]common.PublicKey{},
			amountToSend,
			c.decimals,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not issue additional asset amount: %w", err)
	}
	return tx, nil
}
