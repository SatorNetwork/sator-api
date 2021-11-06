package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) GiveAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, issuer types.Account, recipientAddr string, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))

	asset := common.PublicKeyFromString(assetAddr)
	recipientPublicKey := common.PublicKeyFromString(recipientAddr)
	recipientAta, _, err := common.FindAssociatedTokenAddress(recipientPublicKey, asset)
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
			asset,
			issuer.PublicKey,
			[]common.PublicKey{feePayer.PublicKey},
			amountToSend,
			c.decimals,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not send asset: %w", err)
	}

	return txHash, nil
}

func (c *Client) SendAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, source types.Account, recipientAddr string, amount float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))
	asset := common.PublicKeyFromString(assetAddr)

	sourceAta, _, err := common.FindAssociatedTokenAddress(source.PublicKey, asset)
	if err != nil {
		return "", err
	}

	recipientPublicKey := common.PublicKeyFromString(recipientAddr)
	recipientAta, _, err := common.FindAssociatedTokenAddress(recipientPublicKey, asset)
	if err != nil {
		return "", err
	}

	txHash, err := c.SendTransaction(
		ctx,
		feePayer, source,
		tokenprog.TransferChecked(
			sourceAta,
			recipientAta,
			asset,
			source.PublicKey,
			[]common.PublicKey{feePayer.PublicKey, source.PublicKey},
			amountToSend,
			c.decimals,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not send asset: %w", err)
	}

	return txHash, nil
}
