package puzzle_game

import (
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/client/puzzle_game"
	"github.com/google/uuid"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/utils"
)

func TestTapTile(t *testing.T) {
	defer app_config.RunAndWait()()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	signUpResp, err := c.Auth.SignUp(signUpRequest)
	require.NoError(t, err)
	require.NotNil(t, signUpResp)
	require.NotEmpty(t, signUpResp.AccessToken)

	err = c.Auth.VerifyAcount(signUpResp.AccessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(t, err)

	uid, err := uuid.Parse("7801d5d3-2d2c-4f85-9190-3fa82527f2af")
	require.NoError(t, err)

	_, err = c.PuzzleGame.UnlockPuzzleGame(signUpResp.AccessToken, uid, &puzzle_game.UnlockPuzzleGameRequest{UnlockOption: "899fcefd-9b4b-4c67-905f-b1e580fcaf78"})
	require.NoError(t, err)

	pgBefore, err := c.PuzzleGame.Start(signUpResp.AccessToken, uid)
	require.NoError(t, err)

	pgAfter, err := c.PuzzleGame.TapTile(signUpResp.AccessToken, uid, &puzzle_game.TapTileRequest{
		X: 4,
		Y: 1,
	})
	require.NoError(t, err)

	if !reflect.DeepEqual(pgBefore.Tiles, pgAfter.Tiles) {
		t.Fatalf("no effect from tap tile")
	}
}
