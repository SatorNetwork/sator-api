package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) SendAssets(ctx context.Context, sender types.Account, recipientAddr string, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))
	// Issue Asset
	txHash, err := c.SendTransaction(
		ctx, c.Issuer,
		tokenprog.TransferChecked(
			sender.PublicKey,
			common.PublicKeyFromString(recipientAddr),
			c.Asset.PublicKey,
			c.Issuer.PublicKey,
			[]common.PublicKey{},
			amountToSend,
			c.decimals,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not send Asset: %w", err)
	}

	return txHash, nil
}

func (c *Client) SendAssetsFromIssuer(ctx context.Context, recipientAddr string, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))
	// Issue Asset
	txHash, err := c.SendTransaction(
		ctx, c.Issuer,
		tokenprog.TransferChecked(
			c.Issuer.PublicKey,
			common.PublicKeyFromString(recipientAddr),
			c.Asset.PublicKey,
			c.Issuer.PublicKey,
			[]common.PublicKey{},
			amountToSend,
			c.decimals,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not send Asset: %w", err)
	}

	return txHash, nil
}
