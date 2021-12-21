package two_players

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/internal/test/framework/client"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/internal/test/framework/utils"
)

func TestChallengeV2Sandbox(t *testing.T) {
	t.Skip()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	signUpResp, err := c.Auth.SignUp(signUpRequest)
	require.NoError(t, err)

	defaultChallengeID := "00cfe310-eb1e-42e5-9d6b-801e288a3f8f"
	challenge, err := c.ChallengesClient.GetChallengeById(signUpResp.AccessToken, defaultChallengeID)
	require.NoError(t, err)
	_ = challenge

	// TODO(evg): debug this method
	_, err = c.ChallengesClient.GetQuestionsByChallengeID(signUpResp.AccessToken, defaultChallengeID)
	require.NoError(t, err)
}

func TestDBSandbox(t *testing.T) {
	t.Skip()

	c := client.NewClient()

	ctx := context.Background()
	err := c.DB.Bootstrap(ctx)
	require.NoError(t, err)

	defaultChallengeID, err := c.DB.ChallengeDB().DefaultChallengeID(ctx)
	require.NoError(t, err)
	fmt.Printf("DefaultChallengeID: %v\n", defaultChallengeID)
}
