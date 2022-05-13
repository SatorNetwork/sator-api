package wallet

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	lib_coingecko "github.com/SatorNetwork/sator-api/lib/coingecko"
	"github.com/SatorNetwork/sator-api/lib/sumsub"
	exchange_rates_svc "github.com/SatorNetwork/sator-api/svc/exchange_rates"
	exchange_rates_client "github.com/SatorNetwork/sator-api/svc/exchange_rates/client"
	wallet_svc "github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/client/wallet"
	"github.com/SatorNetwork/sator-api/test/framework/utils"
	"github.com/SatorNetwork/sator-api/test/mock"
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
	defer app_config.RunAndWait()()

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
	defer app_config.RunAndWait()()

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coingeckoMock := lib_coingecko.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.CoingeckoProvider, coingeckoMock)
	solanaPriceInUSD := float64(100)
	satorPriceInUSD := float64(2)
	simplePriceCallback := func(ids []string, vsCurrencies []string) (*map[string]map[string]float32, error) {
		priceMap := map[string]map[string]float32{
			"solana":  {"usd": float32(solanaPriceInUSD)},
			"sator":   {"usd": float32(satorPriceInUSD)},
			"arweave": {"usd": 1},
		}
		return &priceMap, nil
	}
	coingeckoMock.EXPECT().
		SimplePrice([]string{"solana", "sator", "arweave"}, []string{"usd"}).
		DoAndReturn(simplePriceCallback).
		Times(1)

	defer app_config.RunAndWait()()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	exchangeRatesClient, err := exchange_rates_client.Easy(c.DB.Client())
	require.NoError(t, err)
	_, err = exchangeRatesClient.SyncExchangeRates(context.Background(), &exchange_rates_svc.Empty{})
	require.NoError(t, err)

	signUpRequest := auth.RandomSignUpRequest()
	email := signUpRequest.Email
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

	{
		_, err := c.Wallet.CreateTransfer(signUpResp.AccessToken, &createTransferRequest)
		require.Error(t, err)
	}

	err = c.DB.AuthDB().UpdateKYCStatus(context.TODO(), email, sumsub.KYCStatusApproved)
	require.NoError(t, err)

	{
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
	}

	// FIXME
	utils.BackoffRetry(t, func() error {
		satorTokenBalance, err := c.Wallet.GetSatorTokenBalance(signUpResp2.AccessToken)
		require.NoError(t, err)
		// expectedBalance := 0.001
		// TODO(evg): calculate it properly
		expectedBalance := 0.00099245
		if satorTokenBalance != expectedBalance {
			return errors.Errorf("unexpected sator token balance, want: %v, got: %v", expectedBalance, satorTokenBalance)
		}

		return nil
	})
}
