package wallet

import (
	"context"
	
	"github.com/pkg/errors"
	"github.com/portto/solana-go-sdk/common"

	"github.com/SatorNetwork/sator-api/internal/test/framework/accounts"
	wallet_svc "github.com/SatorNetwork/sator-api/svc/wallet"
)

func (w *WalletClient) GetWalletByType(accessToken string, walletType string) (*Wallet, error) {
	wallets, err := w.GetWallets(accessToken)
	if err != nil {
		return nil, errors.Wrap(err, "can't get wallets")
	}
	for _, wallet := range wallets {
		if wallet.Type == walletType {
			return wallet, nil
		}
	}

	return nil, errors.Errorf("%v wallet not found", walletType)
}

func (w *WalletClient) GetSolanaAddress(accessToken string) (string, error) {
	satorWallet, err := w.GetWalletByType(accessToken, wallet_svc.WalletTypeSator)
	if err != nil {
		return "", err
	}
	satorWalletDetails, err := w.GetWalletByID(accessToken, satorWallet.GetDetailsUrl)
	if err != nil {
		return "", err
	}

	return satorWalletDetails.SolanaAccountAddress, nil
}

func (w *WalletClient) GetSolanaPublicKey(accessToken string) (common.PublicKey, error) {
	addr, err := w.GetSolanaAddress(accessToken)
	if err != nil {
		return common.PublicKey{}, err
	}

	return common.PublicKeyFromString(addr), nil
}

func (w *WalletClient) GetSatorTokenPublicKey(accessToken string) (common.PublicKey, error) {
	solanaPublicKey, err := w.GetSolanaPublicKey(accessToken)
	if err != nil {
		return common.PublicKey{}, err
	}
	satorPublicKey, _, err := common.FindAssociatedTokenAddress(solanaPublicKey, accounts.GetAsset().PublicKey)
	if err != nil {
		return common.PublicKey{}, err
	}

	return satorPublicKey, nil
}

func (w *WalletClient) GetSatorTokenAddress(accessToken string) (string, error) {
	addr, err := w.GetSatorTokenPublicKey(accessToken)
	if err != nil {
		return "", err
	}

	return addr.ToBase58(), nil
}

func (w *WalletClient) GetSatorTokenBalance(accessToken string) (float64, error) {
	addr, err := w.GetSatorTokenAddress(accessToken)
	if err != nil {
		return 0, err
	}

	balance, err := w.solanaClient.GetTokenAccountBalance(context.Background(), addr)
	if err != nil {
		return 0, err
	}

	return balance, nil
}
