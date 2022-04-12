//go:build !mock_solana

package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) GiveAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, issuer types.Account, recipientAddr string, amount float64) (string, error) {
	instructions := make([]types.Instruction, 0, 2)
	amountToSend := uint64(amount * float64(c.mltpl))
	asset := common.PublicKeyFromString(assetAddr)

	tokenHolderAta, _, err := common.FindAssociatedTokenAddress(issuer.PublicKey, asset)
	if err != nil {
		return "", err
	}

	recipientPublicKey := common.PublicKeyFromString(recipientAddr)
	recipientAta, err := c.deriveATAPublicKey(ctx, recipientPublicKey, asset)
	if err != nil {
		if !errors.Is(err, ErrATANotCreated) {
			return "", err
		}
		// Add instruction to create token account
		// instructions = append(instructions,
		// 	assotokenprog.CreateAssociatedTokenAccount(
		// 		feePayer.PublicKey,
		// 		recipientPublicKey,
		// 		common.PublicKeyFromString(assetAddr),
		// 	))
		_, err := c.CreateAccountWithATA(ctx, assetAddr, recipientPublicKey.ToBase58(), feePayer)
		if err != nil {
			// return "", fmt.Errorf("CreateAccountWithATA: %w", err)
			log.Printf("CreateAccountWithATA: %v", err)
		}
	}

	instructions = append(instructions, tokenprog.TransferChecked(
		tokenHolderAta,
		recipientAta,
		asset,
		issuer.PublicKey,
		[]common.PublicKey{},
		amountToSend,
		c.decimals,
	))
	txHash, err := c.SendTransaction(ctx, feePayer, issuer, instructions...)
	if err != nil {
		return "", fmt.Errorf("could not send asset: %w", err)
	}

	return txHash, nil
}

func (c *Client) SendAssetsWithAutoDerive(ctx context.Context, assetAddr string, feePayer, source types.Account, recipientAddr, tokenHolderAddr string, amount, fee float64) (string, error) {
	amountToSend := uint64(amount * float64(c.mltpl))
	asset := common.PublicKeyFromString(assetAddr)
	instructions := make([]types.Instruction, 0, 2)

	sourceAta, _, err := common.FindAssociatedTokenAddress(source.PublicKey, asset)
	if err != nil {
		return "", err
	}

	recipientPublicKey := common.PublicKeyFromString(recipientAddr)
	recipientAta, err := c.deriveATAPublicKey(ctx, recipientPublicKey, asset)
	if err != nil {
		if !errors.Is(err, ErrATANotCreated) {
			return "", err
		}
		// Add instruction to create token account
		// instructions = append(instructions,
		// 	assotokenprog.CreateAssociatedTokenAccount(
		// 		feePayer.PublicKey,
		// 		recipientPublicKey,
		// 		common.PublicKeyFromString(assetAddr),
		// 	))

		_, err := c.CreateAccountWithATA(ctx, assetAddr, recipientPublicKey.ToBase58(), feePayer)
		if err != nil {
			// return "", fmt.Errorf("CreateAccountWithATA: %w", err)
			log.Printf("CreateAccountWithATA: %v", err)
		}
	}

	instructions = append(instructions, tokenprog.TransferChecked(
		sourceAta,
		recipientAta,
		asset,
		source.PublicKey,
		[]common.PublicKey{},
		amountToSend,
		c.decimals,
	))

	if fee > 0 {
		tokenHolderPublicKey := common.PublicKeyFromString(tokenHolderAddr)
		tokenHolderAta, err := c.deriveATAPublicKey(ctx, tokenHolderPublicKey, asset)
		if err != nil {
			return "", err
		}

		instructions = append(instructions, tokenprog.TransferChecked(
			sourceAta,
			tokenHolderAta,
			asset,
			source.PublicKey,
			[]common.PublicKey{},
			uint64(fee*float64(c.mltpl)),
			c.decimals,
		))
	}

	txHash, err := c.SendTransaction(ctx, feePayer, source, instructions...)
	if err != nil {
		return "", fmt.Errorf("could not send asset: %w", err)
	}

	return txHash, nil
}
