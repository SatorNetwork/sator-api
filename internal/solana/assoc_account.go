package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/assotokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func (c *Client) CreateAccountWithATA(ctx context.Context, feePayer, issuer, asset, initAcc types.Account) (string, error) {
	instructions := []types.Instruction{
		assotokenprog.CreateAssociatedTokenAccount(
			feePayer.PublicKey,
			initAcc.PublicKey,
			asset.PublicKey,
		),
	}

	res, err := c.solana.GetRecentBlockhash(ctx)
	if err != nil {
		return "", fmt.Errorf("could not get recent block hash: %w", err)
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions:    instructions,
		Signers:         []types.Account{feePayer},
		FeePayer:        feePayer.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return "", fmt.Errorf("could not create new raw transaction: %w", err)
	}

	txhash, err := c.solana.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", fmt.Errorf("could not send raw transaction: %w", err)
	}

	return txhash, nil
}
