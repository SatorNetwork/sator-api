package firebase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	firebase_repository "github.com/SatorNetwork/sator-api/svc/firebase/repository"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestFirebaseDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)

	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()

	defer app_config.RunAndWait()()

	c := client.NewClient()
	ctxb := context.Background()
	firebaseRepository, err := firebase_repository.Prepare(ctxb, c.DB.Client())
	require.NoError(t, err)

	signUpRequest := auth.RandomSignUpRequest()
	user1 := user.NewInitializedUser(signUpRequest, t)
	userID1, err := c.DB.AuthDB().GetUserIDByEmail(context.Background(), user1.Email())
	require.NoError(t, err)
	topic1 := uuid.New().String()

	{
		enabled, err := firebaseRepository.IsNotificationEnabled(ctxb, firebase_repository.IsNotificationEnabledParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)
		require.True(t, enabled)

		err = firebaseRepository.DisableNotification(ctxb, firebase_repository.DisableNotificationParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)

		enabled, err = firebaseRepository.IsNotificationEnabled(ctxb, firebase_repository.IsNotificationEnabledParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)
		require.False(t, enabled)

		err = firebaseRepository.EnableNotification(ctxb, firebase_repository.EnableNotificationParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)

		enabled, err = firebaseRepository.IsNotificationEnabled(ctxb, firebase_repository.IsNotificationEnabledParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)
		require.True(t, enabled)
	}

	{
		disabled, err := firebaseRepository.IsNotificationDisabled(ctxb, firebase_repository.IsNotificationDisabledParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)
		require.False(t, disabled)

		err = firebaseRepository.DisableNotification(ctxb, firebase_repository.DisableNotificationParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)

		disabled, err = firebaseRepository.IsNotificationDisabled(ctxb, firebase_repository.IsNotificationDisabledParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)
		require.True(t, disabled)

		err = firebaseRepository.EnableNotification(ctxb, firebase_repository.EnableNotificationParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)

		disabled, err = firebaseRepository.IsNotificationDisabled(ctxb, firebase_repository.IsNotificationDisabledParams{
			UserID: userID1,
			Topic:  topic1,
		})
		require.NoError(t, err)
		require.False(t, disabled)
	}
}
