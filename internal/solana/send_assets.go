package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) GiveAssetsWithAutoDerive(ctx context.Context, feePayer, issuer, asset types.Account, recipientAddr string, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))

	recipientPublicKey := common.PublicKeyFromString(recipientAddr)
	recipientAta, _, err := common.FindAssociatedTokenAddress(recipientPublicKey, asset.PublicKey)
	if err != nil {
		return "", err
	}

	// Issue asset
	txHash, err := c.SendTransaction(
		ctx,
		feePayer, issuer,
		tokenprog.TransferChecked(
			issuer.PublicKey,
			recipientAta,
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

func (c *Client) SendAssetsWithAutoDerive(ctx context.Context, feePayer, asset, source types.Account, recipientAddr string, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))

	sourceAta, _, err := common.FindAssociatedTokenAddress(source.PublicKey, asset.PublicKey)
	if err != nil {
		return "", err
	}

	recipientPublicKey := common.PublicKeyFromString(recipientAddr)
	recipientAta, _, err := common.FindAssociatedTokenAddress(recipientPublicKey, asset.PublicKey)
	if err != nil {
		return "", err
	}

	txHash, err := c.SendTransaction(
		ctx,
		feePayer, source,
		tokenprog.TransferChecked(
			sourceAta,
			recipientAta,
			asset.PublicKey,
			source.PublicKey,
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
