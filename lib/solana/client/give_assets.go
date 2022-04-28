//go:build !mock_solana

package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) GiveAssetsWithAutoDerive(
	ctx context.Context,
	assetAddr string,
	feePayer types.Account,
	issuer types.Account,
	recipientAddr string,
	amount float64,
) (string, error) {
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

	instructions = append(instructions, tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
		From:     tokenHolderAta,
		To:       recipientAta,
		Mint:     asset,
		Auth:     issuer.PublicKey,
		Signers:  []common.PublicKey{},
		Amount:   amountToSend,
		Decimals: c.decimals,
	}))

	txHash, err := c.SendTransaction(ctx, feePayer, issuer, instructions...)
	if err != nil {
		return "", fmt.Errorf("could not send asset: %w", err)
	}

	return txHash, nil
}
