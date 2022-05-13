package rewards

import (
	"encoding/base64"
	"fmt"
	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
	"testing"
	"time"

	"github.com/SatorNetwork/sator-api/lib/sumsub"
	"github.com/SatorNetwork/sator-api/test/framework/user"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/utils"
)

func TestClaimRewards(t *testing.T) {
	defer app_config.RunAndWait()()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()
	user := user.NewInitializedUser(auth.RandomSignUpRequest(), t)

	err = c.DB.AuthDB().UpdateKYCStatus(context.TODO(), user.Email(), sumsub.KYCStatusApproved)
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
	_ = sc

	userPK, err := c.Wallet.GetSatorTokenPublicKey(user.AccessToken())
	require.NoError(t, err)
	fmt.Println("SATOR user pk:", userPK.ToBase58())

	start := time.Now()
	{
		feePayerPk, err := base64.StdEncoding.DecodeString(app_config.AppConfigForTests.SolanaFeePayerPrivateKey)
		require.NoError(t, err)
		feePayer, err := sc.AccountFromPrivateKeyBytes(feePayerPk)
		require.NoError(t, err)

		//feeAccumulatorPublicKey := common.PublicKeyFromString(app_config.AppConfigForTests.FeeAccumulatorAddress)
		txHash, err := sc.CreateAccountWithATA(context.Background(), app_config.AppConfigForTests.SolanaAssetAddr, userPK.ToBase58(), feePayer)
		require.NoError(t, err)
		fmt.Println(txHash)

		addr, err := c.Wallet.GetSatorTokenAddress(user.AccessToken())
		require.NoError(t, err)
		fmt.Println(addr)
		//balance, err := sc.GetTokenAccountBalance(context.Background(), acc.PublicKey.ToBase58())
		//require.NoError(t, err)
	}
	end := time.Now()
	fmt.Println("CreateAccountWithATA", end.Sub(start))

	err = c.DB.RewardsDB().DepositRewards(context.Background(), id, 100)
	require.NoError(t, err)

	resp, err := c.RewardsClient.ClaimRewards(user.AccessToken())
	require.NoError(t, err)
	require.NotEqual(t, "", resp.TransactionURL)
	fmt.Println(resp.TransactionURL)
}