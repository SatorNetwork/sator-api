package solana

import (
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) SendAssets(feePayer, issuer, asset, sender types.Account, recipientAddr string, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))
	// Issue asset
	txHash, err := c.SendTransaction(
		feePayer, issuer,
		tokenprog.TransferChecked(
			sender.PublicKey,
			common.PublicKeyFromString(recipientAddr),
			asset.PublicKey,
			issuer.PublicKey,
			[]common.PublicKey{},
			amountToSend,
			c.decimals,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not send asset: %w", err)
	}

	return txHash, nil
}
