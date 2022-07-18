package firebase

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	lib_google_firebase "github.com/SatorNetwork/sator-api/lib/google_firebase"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	shows_repository "github.com/SatorNetwork/sator-api/svc/shows/repository"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	shows_client "github.com/SatorNetwork/sator-api/test/framework/client/shows"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestShowsNotifications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)
	messagingMock := lib_google_firebase.NewMockMessagingClientInterface(ctrl)
	mock.RegisterMockObject(mock.GoogleFirebaseMessagingProvider, messagingMock)
	messagingMock.EXPECT().
		Send(gomock.Any(), gomock.Any()).
		Return("", nil).
		Times(2)

	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()

	defer app_config.RunAndWait()()

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	user := user.NewInitializedUser(signUpRequest, t)
	user.SetRole(rbac.RoleAdmin)
	user.RefreshToken()

	show, err := c.ShowsClient.AddShow(user.AccessToken(), &shows_client.AddShowRequest{
		Title:  "title1",
		Cover:  "cover1",
		Status: string(shows_repository.ShowsStatusTypePublished),
	})
	require.NoError(t, err)

	season, err := c.ShowsClient.AddSeason(user.AccessToken(), &shows_client.AddSeasonRequest{
		ShowID:       show.Id,
		SeasonNumber: 1,
	})
	require.NoError(t, err)

	err = c.ShowsClient.AddEpisode(user.AccessToken(), &shows_client.AddEpisodeRequest{
		ShowID:      show.Id,
		SeasonID:    season.Id,
		Title:       "title1",
		ReleaseDate: time.Now().Format(time.RFC3339),
		Status:      string(shows_repository.EpisodesStatusTypePublished),
	})
	require.NoError(t, err)
}
