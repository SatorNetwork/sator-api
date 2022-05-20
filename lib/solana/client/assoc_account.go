//go:build !mock_solana

package client

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/assotokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) CreateAccountWithATA(ctx context.Context, assetAddr, initAccAddr string, feePayer types.Account) (string, error) {
	asset := common.PublicKeyFromString(assetAddr)
	initAcc := common.PublicKeyFromString(initAccAddr)

	initAccAta, _, err := c.FindAssociatedTokenAddress(initAcc, asset)
	if err != nil {
		return "", errors.Wrap(err, "can't find ata error")
	}

	instructions := []types.Instruction{
		assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
			Funder:                 feePayer.PublicKey,
			Owner:                  initAcc,
			Mint:                   asset,
			AssociatedTokenAccount: initAccAta,
		}),
	}

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			Instructions:    instructions,
			RecentBlockhash: res.Blockhash,
		}),
		Signers: []types.Account{feePayer},
	})
	if err != nil {
		return "", fmt.Errorf("could not create new raw transaction: %w", err)
	}

	txhash, ok, err := c.SendTransactionUntilConfirmed(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("could not send transaction: %w", err)
	}
	if !ok {
		return "", errors.New("tx not confirmed")
	}

	recipientAta, err := c.deriveATAPublicKey(ctx, initAcc, asset)
	if err != nil {
		return "", err
	}
	_ = recipientAta

	return txhash, nil
}
