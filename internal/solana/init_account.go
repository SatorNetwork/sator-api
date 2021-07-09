package solana

import (
	"context"
	"fmt"

	"github.com/portto/solana-go-sdk/types"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/sysprog"
	"github.com/portto/solana-go-sdk/tokenprog"
)

func (c *Client) InitAccountToUseAsset(ctx context.Context, initAcc types.Account) (string, error) {
	// Allow account hold the new Asset
	rentExemptionBalanceIssuer, err := c.solana.GetMinimumBalanceForRentExemption(ctx, tokenprog.TokenAccountSize)
	if err != nil {
		return "", fmt.Errorf("could not get min balance for rent exemption: %w", err)
	}

	tx, err := c.SendTransaction(
		ctx, initAcc,
		sysprog.CreateAccount(
			c.FeePayer.PublicKey,
			initAcc.PublicKey,
			common.TokenProgramID,
			rentExemptionBalanceIssuer,
			tokenprog.TokenAccountSize,
		),
		tokenprog.InitializeAccount(
			initAcc.PublicKey,
			c.Asset.PublicKey,
			c.Issuer.PublicKey,
		),
	)
	if err != nil {
		return "", fmt.Errorf("could not associate account: %w", err)
	}

	return tx, nil
}
