package rewards

import (
	"encoding/base64"
	"fmt"
	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
	"github.com/SatorNetwork/sator-api/lib/sumsub"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/mock"
	"github.com/golang/mock/gomock"
	solana_sdk_client "github.com/portto/solana-go-sdk/client"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestClaimRewards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)

	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.EXPECT().
		GetTokenAccountBalance(gomock.Any(), gomock.Any()).
		Return(float64(0), nil).
		AnyTimes()
	solanaMock.EXPECT().
		RequestAirdrop(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).
		AnyTimes()
	solanaMock.EXPECT().
		GetAccountBalanceSOL(gomock.Any(), gomock.Any()).
		Return(float64(0), nil).
		AnyTimes()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()
	solanaMock.EXPECT().
		CreateAccountWithATA(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).
		AnyTimes()
	solanaMock.EXPECT().
		GetTokenAccountBalanceWithAutoDerive(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(float64(101), nil).
		AnyTimes()
	solanaMock.EXPECT().
		SendAssetsWithAutoDerive(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).
		AnyTimes()

	cnt := 1
	getTransactionsCallback := func(ctx context.Context, txhash string) (*solana_sdk_client.GetTransactionResponse, error) {
		defer func() {
			cnt++
		}()

		fmt.Println(cnt)

		if cnt <= 2 {
			return nil, fmt.Errorf("err")
		}

		return &solana_sdk_client.GetTransactionResponse{}, nil
	}

	solanaMock.EXPECT().
		GetTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(getTransactionsCallback).
		Times(3)

	defer app_config.RunAndWait()()

	c := client.NewClient()
	user := user.NewInitializedUser(auth.RandomSignUpRequest(), t)

	err := c.DB.AuthDB().UpdateKYCStatus(context.TODO(), user.Email(), sumsub.KYCStatusApproved)
	require.NoError(t, err)

	id, err := c.DB.AuthDB().GetUserIDByEmail(context.Background(), user.Email())
	require.NoError(t, err)

	sc := solana_client.New(app_config.AppConfigForTests.SolanaApiBaseUrl, solana_client.Config{
		SystemProgram:  app_config.AppConfigForTests.SolanaSystemProgram,
		SysvarRent:     app_config.AppConfigForTests.SolanaSysvarRent,
		SysvarClock:    app_config.AppConfigForTests.SolanaSysvarClock,
		SplToken:       app_config.AppConfigForTests.SolanaSplToken,
		StakeProgramID: app_config.AppConfigForTests.SolanaStakeProgramID,
	}, nil)

	userPK, err := c.Wallet.GetSatorTokenPublicKey(user.AccessToken())
	require.NoError(t, err)

	//var tokenAccount string
	{
		feePayerPk, err := base64.StdEncoding.DecodeString(app_config.AppConfigForTests.SolanaFeePayerPrivateKey)
		require.NoError(t, err)
		feePayer, err := sc.AccountFromPrivateKeyBytes(feePayerPk)
		require.NoError(t, err)

		txHash, err := sc.CreateAccountWithATA(context.Background(), app_config.AppConfigForTests.SolanaAssetAddr, userPK.ToBase58(), feePayer)
		require.NoError(t, err)
		_ = txHash

		_, err = c.Wallet.GetSatorTokenAddress(user.AccessToken())
		require.NoError(t, err)
	}

	err = c.DB.RewardsDB().DepositRewards(context.Background(), id, 100)
	require.NoError(t, err)

	resp, err := c.RewardsClient.ClaimRewards(user.AccessToken())
	require.NoError(t, err)
	require.NotEqual(t, "", resp.TransactionURL)

	time.Sleep(time.Minute * 3)
}
