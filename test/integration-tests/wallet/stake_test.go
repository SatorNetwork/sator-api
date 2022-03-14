package wallet

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	wallet_svc "github.com/SatorNetwork/sator-api/svc/wallet"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	wallet_client "github.com/SatorNetwork/sator-api/test/framework/client/wallet"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestStake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)
	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()
	solanaMock.EXPECT().
		GetTokenAccountBalanceWithAutoDerive(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(100.0, nil).
		AnyTimes()
	solanaMock.EXPECT().
		Stake(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).
		AnyTimes()
	solanaMock.EXPECT().
		Unstake(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).
		AnyTimes()

	defer app_config.RunAndWait()()

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

	wallet, err := c.Wallet.GetWalletByType(signUpResp.AccessToken, wallet_svc.WalletTypeSator)
	require.NoError(t, err)

	_, err = c.Wallet.GetStake(signUpResp.AccessToken, wallet.Id)
	require.NoError(t, err)
	err = c.Wallet.SetStake(signUpResp.AccessToken, &wallet_client.SetStakeRequest{
		Amount:   10,
		WalletID: wallet.Id,
		Duration: 0,
	})
	require.NoError(t, err)
	_, err = c.Wallet.GetStake(signUpResp.AccessToken, wallet.Id)
	require.NoError(t, err)
	time.Sleep(time.Second)
	err = c.Wallet.Unstake(signUpResp.AccessToken, wallet.Id)
	require.NoError(t, err)
}
