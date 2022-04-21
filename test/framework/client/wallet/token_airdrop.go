package wallet

import (
	"context"

	"github.com/SatorNetwork/sator-api/test/framework/accounts"
)

func (w *WalletClient) RequestTokenAirdrop(solanaAddr string, amount float64) error {
	feePayer, tokenHolder, asset := accounts.GetAccounts()

	_, err := w.solanaClient.GiveAssetsWithAutoDerive(
		context.Background(),
		asset.PublicKey.ToBase58(),
		feePayer,
		tokenHolder,
		solanaAddr,
		amount,
	)
	return err
}
