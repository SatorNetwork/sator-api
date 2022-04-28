package shows

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/portto/solana-go-sdk/types"
	"github.com/stretchr/testify/require"

	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	shows_repository "github.com/SatorNetwork/sator-api/svc/shows/repository"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/client/shows"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestSendTips(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)
	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()
	solanaMock.EXPECT().
		GetTokenAccountBalanceWithAutoDerive(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(100.0, nil).
		AnyTimes()
	sendAssetsCallback := func(
		ctx context.Context,
		assetAddr string,
		feePayer types.Account,
		source types.Account,
		recipientAddr string,
		amount float64,
		cfg *solana_lib.SendAssetsConfig,
	) (string, error) {
		require.Equal(t, float64(1), amount)
		require.Equal(t, app_config.AppConfigForTests.TipsPercent, cfg.PercentToCharge)
		require.Equal(t, true, cfg.ChargeSolanaFeeFromSender)
		return "", nil
	}
	solanaMock.EXPECT().
		SendAssetsWithAutoDerive(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(sendAssetsCallback).
		AnyTimes()

	defer app_config.RunAndWait()()

	c := client.NewClient()
	defaultChallengeID, err := c.DB.ChallengeDB().DefaultChallengeID(context.Background())
	require.NoError(t, err)

	user := user.NewInitializedUser(auth.RandomSignUpRequest(), t)
	userID, err := c.DB.AuthDB().GetUserIDByEmail(context.Background(), user.Email())
	require.NoError(t, err)

	episodeID, err := c.DB.ShowsDB().Repository().GetEpisodeIDByQuizChallengeID(context.Background(), uuid.NullUUID{
		UUID:  defaultChallengeID,
		Valid: true,
	})
	require.NoError(t, err)

	rating, err := c.DB.ShowsDB().Repository().ReviewEpisode(context.Background(), shows_repository.ReviewEpisodeParams{
		EpisodeID: episodeID,
		UserID:    userID,
		Rating:    1,
	})
	require.NoError(t, err)

	err = c.ShowsClient.SendTipsToReviewAuthor(user.AccessToken(), &shows.SendTipsRequest{
		ReviewID: rating.ID.String(),
		Amount:   1,
	})
	require.NoError(t, err)
}
