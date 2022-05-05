package rewards

import (
	"fmt"
	"testing"

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

	err = c.DB.RewardsDB().DepositRewards(context.Background(), id, 100)
	require.NoError(t, err)

	resp, err := c.RewardsClient.ClaimRewards(user.AccessToken())
	require.NoError(t, err)
	require.NotEqual(t, "", resp.TransactionURL)
	fmt.Println(resp.TransactionURL)
}