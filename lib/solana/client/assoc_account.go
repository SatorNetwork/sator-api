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
			Owner:                  common.PublicKeyFromString(initAccAddr),
			Mint:                   common.PublicKeyFromString(assetAddr),
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

	txhash, err := c.solana.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("could not send transaction: %w", err)
	}

	return txhash, nil
}
