package settings

import (
	"context"
	"github.com/SatorNetwork/sator-api/test/framework/client/settings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/framework/utils"
)

func TestAddSetting(t *testing.T) {
	defer app_config.RunAndWait()()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	u := user.NewInitializedUser(signUpRequest, t)
	u.SetRole(rbac.RoleAdmin)
	u.RefreshToken()

	req := &settings.Setting{
		Key:         "test",
		Name:        "test",
		Value:       true,
		ValueType:   "bool",
		Description: "",
	}
	setting, err := c.SettingsClient.AddSetting(u.AccessToken(), req)
	require.NoError(t, err)
	require.Equal(t, setting, req)
}

func TestUpdateSetting(t *testing.T) {
	defer app_config.RunAndWait()()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	u := user.NewInitializedUser(signUpRequest, t)
	u.SetRole(rbac.RoleAdmin)
	u.RefreshToken()

	req := &settings.Setting{
		Key:         "test",
		Name:        "test",
		Value:       true,
		ValueType:   "bool",
		Description: "",
	}
	setting, err := c.SettingsClient.AddSetting(u.AccessToken(), req)
	require.NoError(t, err)

	setting.Name = "test123"

	setting, err = c.SettingsClient.UpdateSetting(u.AccessToken(), setting)
	require.NoError(t, err)
	require.Equal(t, "test123", setting.Name)
}
