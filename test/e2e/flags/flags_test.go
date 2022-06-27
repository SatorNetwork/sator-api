package flags

import (
	"context"
	"github.com/SatorNetwork/sator-api/svc/flags/alias"
	"github.com/SatorNetwork/sator-api/test/framework/client/flags"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/framework/utils"
)

func TestGetFlags(t *testing.T) {
	t.Skip()
	defer app_config.RunAndWait()()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	u := user.NewInitializedUser(signUpRequest, t)
	u.SetRole(rbac.RoleAdmin)

	flagsList, err := c.FlagsClient.GetFlags(u.AccessToken())
	require.NoError(t, err)
	require.NotEmpty(t, flagsList)
	require.GreaterOrEqual(t, len(flagsList), 1)
}

func TestUpdateFlag(t *testing.T) {
	t.Skip()
	defer app_config.RunAndWait()()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	u := user.NewInitializedUser(signUpRequest, t)
	u.SetRole(rbac.RoleAdmin)
	u.SignUp()

	flag, err := c.FlagsClient.UpdateFlag(u.AccessToken(), &flags.Flag{
		Key:   alias.FlagKeyPuzzleGameRewards.String(),
		Value: alias.FlagValueDisabled.String(),
	})
	require.NoError(t, err)
	require.Equal(t, alias.FlagValueDisabled.String(), flag.Value)
}
