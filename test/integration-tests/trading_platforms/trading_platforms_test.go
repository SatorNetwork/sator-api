package trading_platforms

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/lib/rbac"
	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	trading_platforms_client "github.com/SatorNetwork/sator-api/test/framework/client/trading_platforms"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestTradingPlatform(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)
	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()

	defer app_config.RunAndWait()()

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	user := user.NewInitializedUser(signUpRequest, t)
	user.SetRole(rbac.RoleAdmin)
	user.RefreshToken()

	var linkID string
	{
		resp, err := c.TradingPlatformsClient.CreateLink(user.AccessToken(), &trading_platforms_client.CreateLinkRequest{
			Title: "title",
			Link:  "link",
			Logo:  "logo",
		})
		t.Logf("c.TradingPlatformsClient.CreateLink: %v", err)
		require.NoError(t, err)
		linkID = resp.Id

		links, err := c.TradingPlatformsClient.GetLinks(user.AccessToken(), &trading_platforms_client.Empty{})
		require.NoError(t, err)
		linkIDs := getLinkIDs(links)
		require.Contains(t, linkIDs, linkID)

		link, err := getLinkByID(links, linkID)
		require.NoError(t, err)
		require.Equal(t, "title", link.Title)
		require.Equal(t, "link", link.Link)
		require.Equal(t, "logo", link.Logo)
	}

	{
		_, err := c.TradingPlatformsClient.UpdateLink(user.AccessToken(), linkID, &trading_platforms_client.UpdateLinkRequest{
			Title: "title-updated",
			Link:  "link-updated",
			Logo:  "logo-updated",
		})
		require.NoError(t, err)

		links, err := c.TradingPlatformsClient.GetLinks(user.AccessToken(), &trading_platforms_client.Empty{})
		require.NoError(t, err)
		linkIDs := getLinkIDs(links)
		require.Contains(t, linkIDs, linkID)

		link, err := getLinkByID(links, linkID)
		require.NoError(t, err)
		require.Equal(t, "title-updated", link.Title)
		require.Equal(t, "link-updated", link.Link)
		require.Equal(t, "logo-updated", link.Logo)
	}

	{
		_, err := c.TradingPlatformsClient.DeleteLink(user.AccessToken(), linkID, &trading_platforms_client.Empty{})
		require.NoError(t, err)

		links, err := c.TradingPlatformsClient.GetLinks(user.AccessToken(), &trading_platforms_client.Empty{})
		require.NoError(t, err)
		linkIDs := getLinkIDs(links)
		require.NotContains(t, linkIDs, linkID)
	}
}

func getLinkIDs(links []*trading_platforms_client.Link) []string {
	linkIDs := make([]string, 0, len(links))
	for _, link := range links {
		linkIDs = append(linkIDs, link.Id)
	}

	return linkIDs
}

func getLinkByID(links []*trading_platforms_client.Link, id string) (*trading_platforms_client.Link, error) {
	for _, link := range links {
		if link.Id == id {
			return link, nil
		}
	}

	return nil, errors.Errorf("can't get link by id: %v", id)
}
