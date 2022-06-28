package settings

import (
	"context"
	"github.com/SatorNetwork/sator-api/test/framework/client/settings"
	"github.com/google/uuid"
	"strings"
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

	uid, err := uuid.NewUUID()
	require.NoError(t, err)

	key := toSnakeCase("test-" + uid.String())
	req := &settings.Setting{
		Key:         key,
		Name:        "Test",
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

	uid, err := uuid.NewUUID()
	require.NoError(t, err)

	key := "test" + uid.String()
	req := &settings.Setting{
		Key:         key,
		Name:        "Test",
		Value:       true,
		ValueType:   "bool",
		Description: "",
	}
	setting, err := c.SettingsClient.AddSetting(u.AccessToken(), req)
	require.NoError(t, err)

	setting.Value = false

	setting, err = c.SettingsClient.UpdateSetting(u.AccessToken(), setting)
	require.NoError(t, err)
	require.Equal(t, false, setting.Value)
}

// format string to snake case
func toSnakeCase(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = alphaNumUnderscore(s)

	return s
}

// get only alphanumeric characters from string and replace spaces with underscores
func alphaNumUnderscore(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			return r
		}
		if r == ' ' || r == '_' || r == '-' || r == '.' || r == ',' || r == ':' || r == ';' || r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' {
			return '_'
		}
		return -1
	}, s)
}
