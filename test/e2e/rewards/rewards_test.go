package rewards

import (
	"encoding/base64"
	"testing"
	"time"

	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"

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

	userPK, err := c.Wallet.GetSatorTokenPublicKey(user.AccessToken())
	require.NoError(t, err)

	var tokenAccount string
	{
		feePayerPk, err := base64.StdEncoding.DecodeString(app_config.AppConfigForTests.SolanaFeePayerPrivateKey)
		require.NoError(t, err)
		feePayer, err := sc.AccountFromPrivateKeyBytes(feePayerPk)
		require.NoError(t, err)

		txHash, err := sc.CreateAccountWithATA(context.Background(), app_config.AppConfigForTests.SolanaAssetAddr, userPK.ToBase58(), feePayer)
		require.NoError(t, err)
		_ = txHash

		tokenAccount, err = c.Wallet.GetSatorTokenAddress(user.AccessToken())
		require.NoError(t, err)
	}

	err = c.DB.RewardsDB().DepositRewards(context.Background(), id, 100)
	require.NoError(t, err)

	resp, err := c.RewardsClient.ClaimRewards(user.AccessToken())
	require.NoError(t, err)
	require.NotEqual(t, "", resp.TransactionURL)

	time.Sleep(time.Second * 65)

	var balance float64
	i := 0
	ticker := time.NewTicker(time.Minute)
LOOP:
	for {
		select {
		case <-ticker.C:
			balance, err = sc.GetTokenAccountBalance(context.Background(), tokenAccount)
			require.NoError(t, err)
			if balance == 99.25 {
				break LOOP
			}
			if i == 20 {
				t.Fatalf("transaction timeout")
			}
			i++
		}
	}
}
