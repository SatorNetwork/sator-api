package wallet

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/internal/test/framework/client"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/wallet"
	"github.com/SatorNetwork/sator-api/internal/test/framework/utils"
	wallet_svc "github.com/SatorNetwork/sator-api/svc/wallet"
)

func isWalletValid(w *wallet.Wallet) bool {
	return w.Id != "" &&
		w.Type != "" &&
		w.GetDetailsUrl != "" &&
		w.GetTransactionsUrl != "" &&
		w.Order != 0
}

func isWalletDetailsValid(w *wallet.WalletDetails) bool {
	return w.Id != "" &&
		w.Order != 0 &&
		w.Balance != nil &&
		w.Actions != nil
}

func isTxValid(tx *wallet.Tx) bool {
	return tx.Id != "" &&
		tx.WalletId != "" &&
		tx.TxHash != "" &&
		tx.Amount != 0 &&
		tx.CreatedAt != ""
}

func isCreateTransferResponseValid(resp *wallet.CreateTransferResponse) bool {
	return resp.AssetName != "" &&
		resp.Amount != 0 &&
		resp.RecipientAddress != "" &&
		resp.TxHash != "" &&
		resp.SenderWalletId != ""
}

func TestWalletCreation(t *testing.T) {
	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	signUpResp, err := c.Auth.SignUp(signUpRequest)
	require.NoError(t, err)
	require.NotNil(t, signUpResp)
	require.NotEmpty(t, signUpResp.AccessToken)

	err = c.Auth.VerifyAcount(signUpResp.AccessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(t, err)

	// check GetWallets API method
	wallets, err := c.Wallet.GetWallets(signUpResp.AccessToken)
	require.NoError(t, err)
	require.NotNil(t, wallets)
	walletNum := 2
	require.Len(t, wallets, walletNum)
	for _, w := range wallets {
		require.NotNil(t, w)
		require.True(t, isWalletValid(w))
	}

	// check GetWalletByID API method
	{
		for _, w := range wallets {
			walletDetails, err := c.Wallet.GetWalletByID(signUpResp.AccessToken, w.GetDetailsUrl)
			require.NoError(t, err)
			require.NotNil(t, walletDetails)
			require.True(t, isWalletDetailsValid(walletDetails))
		}
	}

	// check GetWalletByType helper
	{
		satorWallet, err := c.Wallet.GetWalletByType(signUpResp.AccessToken, wallet_svc.WalletTypeSator)
		require.NoError(t, err)
		require.NotNil(t, satorWallet)
		require.True(t, isWalletValid(satorWallet))

		rewardsWallet, err := c.Wallet.GetWalletByType(signUpResp.AccessToken, wallet_svc.WalletTypeRewards)
		require.NoError(t, err)
		require.NotNil(t, rewardsWallet)
		require.True(t, isWalletValid(rewardsWallet))
	}
}

func TestGetWalletTxsAPI(t *testing.T) {
	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	signUpResp, err := c.Auth.SignUp(signUpRequest)
	require.NoError(t, err)
	require.NotNil(t, signUpResp)
	require.NotEmpty(t, signUpResp.AccessToken)

	err = c.Auth.VerifyAcount(signUpResp.AccessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(t, err)

	wallets, err := c.Wallet.GetWallets(signUpResp.AccessToken)
	require.NoError(t, err)

	for _, w := range wallets {
		txs, err := c.Wallet.GetWalletTxs(signUpResp.AccessToken, w.GetTransactionsUrl)
		require.NoError(t, err)
		require.NotNil(t, txs)
		for _, tx := range txs {
			require.NotNil(t, tx)
			require.True(t, isTxValid(tx))
		}
	}
}

func TestSPLTokenPayment(t *testing.T) {
	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	signUpResp, err := c.Auth.SignUp(signUpRequest)
	require.NoError(t, err)
	require.NotNil(t, signUpResp)
	require.NotEmpty(t, signUpResp.AccessToken)

	err = c.Auth.VerifyAcount(signUpResp.AccessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(t, err)

	signUpRequest = auth.RandomSignUpRequest()
	signUpResp2, err := c.Auth.SignUp(signUpRequest)
	require.NoError(t, err)
	require.NotNil(t, signUpResp2)
	require.NotEmpty(t, signUpResp2.AccessToken)

	err = c.Auth.VerifyAcount(signUpResp2.AccessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(t, err)

	satorWallet, err := c.Wallet.GetWalletByType(signUpResp.AccessToken, wallet_svc.WalletTypeSator)
	require.NoError(t, err)
	solanaAddr, err := c.Wallet.GetSolanaAddress(signUpResp.AccessToken)
	require.NoError(t, err)
	solanaAddr2, err := c.Wallet.GetSolanaAddress(signUpResp2.AccessToken)
	require.NoError(t, err)

	{
		utils.BackoffRetry(t, func() error {
			err = c.Wallet.RequestTokenAirdrop(solanaAddr, 1)
			return err
		})

		utils.BackoffRetry(t, func() error {
			satorTokenBalance, err := c.Wallet.GetSatorTokenBalance(signUpResp.AccessToken)
			require.NoError(t, err)
			if satorTokenBalance != 1 {
				return errors.Errorf("unexpected sator token balance, want: %v, got: %v", 1, satorTokenBalance)
			}

			return nil
		})
	}

	createTransferRequest := wallet_svc.CreateTransferRequest{
		SenderWalletID:   satorWallet.Id,
		RecipientAddress: solanaAddr2,
		Amount:           0.001,
		Asset:            "",
	}
	resp, err := c.Wallet.CreateTransfer(signUpResp.AccessToken, &createTransferRequest)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.True(t, isCreateTransferResponseValid(resp))
	confirmTransferRequest := wallet_svc.ConfirmTransferRequest{
		SenderWalletID:  satorWallet.Id,
		TransactionHash: resp.TxHash,
	}
	err = c.Wallet.ConfirmTransfer(signUpResp.AccessToken, &confirmTransferRequest)
	require.NoError(t, err)

	utils.BackoffRetry(t, func() error {
		satorTokenBalance, err := c.Wallet.GetSatorTokenBalance(signUpResp2.AccessToken)
		require.NoError(t, err)
		if satorTokenBalance != 0.001 {
			return errors.Errorf("unexpected sator token balance, want: %v, got: %v", 0.001, satorTokenBalance)
		}

		return nil
	})
}
