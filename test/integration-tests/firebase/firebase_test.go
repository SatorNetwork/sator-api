package firebase

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	solana_lib "github.com/SatorNetwork/sator-api/lib/solana"
	firebase_repository "github.com/SatorNetwork/sator-api/svc/firebase/repository"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/client/firebase"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/mock"
)

func TestFirebase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	solanaMock := solana_lib.NewMockInterface(ctrl)
	mock.RegisterMockObject(mock.SolanaProvider, solanaMock)

	solanaMock.ExpectCheckPrivateKeyAny()
	solanaMock.ExpectNewAccountAny()
	solanaMock.ExpectAccountFromPrivateKeyBytesAny()

	defer app_config.RunAndWait()()

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	user := user.NewInitializedUser(signUpRequest, t)

	deviceID := "test-device-id"

	{
		token := "test-token"
		_, err := c.FirebaseClient.RegisterToken(user.AccessToken(), &firebase.RegisterTokenRequest{
			DeviceId: deviceID,
			Token:    token,
		})
		require.NoError(t, err)

		userID, err := c.DB.AuthDB().GetUserIDByEmail(context.Background(), user.Email())
		require.NoError(t, err)
		firebaseRepo := c.DB.FirebaseDB().Repository()
		resp, err := firebaseRepo.GetRegistrationToken(context.Background(), firebase_repository.GetRegistrationTokenParams{
			DeviceID: deviceID,
			UserID:   userID,
		})
		require.NoError(t, err)
		require.Equal(t, deviceID, resp.DeviceID)
		require.Equal(t, userID, resp.UserID)
		require.Equal(t, token, resp.RegistrationToken)
	}

	{
		token2 := "test-token-2"
		_, err := c.FirebaseClient.RegisterToken(user.AccessToken(), &firebase.RegisterTokenRequest{
			DeviceId: deviceID,
			Token:    token2,
		})
		require.NoError(t, err)

		userID, err := c.DB.AuthDB().GetUserIDByEmail(context.Background(), user.Email())
		require.NoError(t, err)
		firebaseRepo := c.DB.FirebaseDB().Repository()
		resp, err := firebaseRepo.GetRegistrationToken(context.Background(), firebase_repository.GetRegistrationTokenParams{
			DeviceID: deviceID,
			UserID:   userID,
		})
		require.NoError(t, err)
		require.Equal(t, deviceID, resp.DeviceID)
		require.Equal(t, userID, resp.UserID)
		require.Equal(t, token2, resp.RegistrationToken)
	}
}
