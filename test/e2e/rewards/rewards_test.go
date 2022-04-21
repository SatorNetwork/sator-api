package rewards

import (
	"github.com/SatorNetwork/sator-api/lib/sumsub"
	"testing"

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

	signUpRequest := auth.RandomSignUpRequest()
	signUpResp, err := c.Auth.SignUp(signUpRequest)
	require.NoError(t, err)
	require.NotNil(t, signUpResp)
	require.NotEmpty(t, signUpResp.AccessToken)

	err = c.DB.AuthDB().UpdateKYCStatus(context.TODO(), signUpRequest.Email, sumsub.KYCStatusApproved)
	require.NoError(t, err)

	id, err := c.DB.AuthDB().GetUserIDByEmail(context.Background(), signUpRequest.Email)
	require.NoError(t, err)

	err = c.Auth.VerifyAcount(signUpResp.AccessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(t, err)

	err = c.DB.RewardsDB().DepositRewards(context.Background(), id, 100)
	require.NoError(t, err)

	resp, err := c.RewardsClient.ClaimRewards(signUpResp.AccessToken)
	require.NoError(t, err)
	require.NotEqual(t, resp.TransactionURL, "")

	_, err = c.RewardsClient.ClaimRewards(signUpResp.AccessToken)
	require.Error(t, err)
}
