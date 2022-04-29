package puzzle_game

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	pgsvc "github.com/SatorNetwork/sator-api/svc/puzzle_game"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/client/puzzle_game"
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

	shows, err := c.DB.ShowsDB().GetShowsByTitle(context.Background(), "test-title-show")
	require.NoError(t, err)

	if len(shows) != 1 {
		t.Fatalf("show len is not 1, len : %v", len(shows))
	}

	episodes, err := c.DB.ShowsDB().GetEpisodesIDByShowID(shows[0].ID)
	require.NoError(t, err)

	if len(episodes) != 1 {
		t.Fatalf("episodes len is not 1, len : %v", len(episodes))
	}

	pg, err := c.DB.PuzzleGameDB().GetPuzzleGameByEpisodeID(context.Background(), episodes[0])
	require.NoError(t, err)

	upg, err := c.PuzzleGameClient.UnlockPuzzleGame(signUpResp.AccessToken, pg.ID, &puzzle_game.UnlockPuzzleGameRequest{UnlockOption: "899fcefd-9b4b-4c67-905f-b1e580fcaf78"})
	require.NoError(t, err)
	_ = upg

	pgBefore, err := c.PuzzleGameClient.Start(signUpResp.AccessToken, pg.ID)
	require.NoError(t, err)

	pgAfter, err := c.PuzzleGameClient.TapTile(signUpResp.AccessToken, pg.ID, &puzzle_game.TapTileRequest{
		X: 4,
		Y: 1,
	})
	require.NoError(t, err)
	_ = pgBefore

	for _, tile := range pgAfter.Tiles {
		if tile.IsWhitespace {
			require.Equal(t, tile.CurrentPosition.X, 4)
			require.Equal(t, tile.CurrentPosition.Y, 1)
		}
	}

	pgBefore = pgAfter
	pgAfter, err = c.PuzzleGameClient.TapTile(signUpResp.AccessToken, pg.ID, &puzzle_game.TapTileRequest{
		X: 1,
		Y: 1,
	})
	require.NoError(t, err)

	for _, tile := range pgAfter.Tiles {
		if tile.IsWhitespace {
			require.Equal(t, tile.CurrentPosition.X, 1)
			require.Equal(t, tile.CurrentPosition.Y, 1)
		}
	}

	pgAfter, err = c.PuzzleGameClient.TapTile(signUpResp.AccessToken, pg.ID, &puzzle_game.TapTileRequest{
		X: 4,
		Y: 1,
	})
	require.NoError(t, err)

	pgAfter, err = c.PuzzleGameClient.TapTile(signUpResp.AccessToken, pg.ID, &puzzle_game.TapTileRequest{
		X: 4,
		Y: 4,
	})
	require.NoError(t, err)
	require.Equal(t, pgsvc.PuzzleGameStatusFinished, int(pgAfter.Status))
}

func TestTapTileStepLimit(t *testing.T) {
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

	shows, err := c.DB.ShowsDB().GetShowsByTitle(context.Background(), "test-title-show")
	require.NoError(t, err)

	if len(shows) != 1 {
		t.Fatalf("show len is not 1, len : %v", len(shows))
	}

	episodes, err := c.DB.ShowsDB().GetEpisodesIDByShowID(shows[0].ID)
	require.NoError(t, err)

	if len(episodes) != 1 {
		t.Fatalf("episodes len is not 1, len : %v", len(episodes))
	}

	pg, err := c.DB.PuzzleGameDB().GetPuzzleGameByEpisodeID(context.Background(), episodes[0])
	require.NoError(t, err)

	upg, err := c.PuzzleGameClient.UnlockPuzzleGame(signUpResp.AccessToken, pg.ID, &puzzle_game.UnlockPuzzleGameRequest{UnlockOption: "899fcefd-9b4b-4c67-905f-b1e580fcaf78"})
	require.NoError(t, err)
	_ = upg

	pgBefore, err := c.PuzzleGameClient.Start(signUpResp.AccessToken, pg.ID)
	require.NoError(t, err)

	pgAfter, err := c.PuzzleGameClient.TapTile(signUpResp.AccessToken, pg.ID, &puzzle_game.TapTileRequest{
		X: 4,
		Y: 1,
	})
	require.NoError(t, err)
	_ = pgBefore

	for _, tile := range pgAfter.Tiles {
		if tile.IsWhitespace {
			require.Equal(t, tile.CurrentPosition.X, 4)
			require.Equal(t, tile.CurrentPosition.Y, 1)
		}
	}

	for i := 0; i < 15; i++ {
		pgAfter, err = c.PuzzleGameClient.TapTile(signUpResp.AccessToken, pg.ID, &puzzle_game.TapTileRequest{
			X: 1,
			Y: 1,
		})
		require.NoError(t, err)

		for _, tile := range pgAfter.Tiles {
			if tile.IsWhitespace {
				require.Equal(t, tile.CurrentPosition.X, 1)
				require.Equal(t, tile.CurrentPosition.Y, 1)
			}
		}

		pgAfter, err = c.PuzzleGameClient.TapTile(signUpResp.AccessToken, pg.ID, &puzzle_game.TapTileRequest{
			X: 4,
			Y: 1,
		})
		require.NoError(t, err)

		for _, tile := range pgAfter.Tiles {
			if tile.IsWhitespace {
				require.Equal(t, tile.CurrentPosition.X, 4)
				require.Equal(t, tile.CurrentPosition.Y, 1)
			}
		}
	}


	pgAfter, err = c.PuzzleGameClient.TapTile(signUpResp.AccessToken, pg.ID, &puzzle_game.TapTileRequest{
		X: 1,
		Y: 1,
	})
	require.NoError(t, err)
	require.Equal(t, pgsvc.PuzzleStatusReachedStepLimit, int(pgAfter.Status))

	pgAfter, err = c.PuzzleGameClient.TapTile(signUpResp.AccessToken, pg.ID, &puzzle_game.TapTileRequest{
		X: 4,
		Y: 1,
	})
	require.Error(t, err)
}