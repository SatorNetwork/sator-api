package ten_players

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/lib/encryption/envelope"
	internal_rsa "github.com/SatorNetwork/sator-api/lib/encryption/rsa"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/e2e/quiz_v2/message_container"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/message_verifier"
	"github.com/SatorNetwork/sator-api/test/framework/nats_subscriber"
	"github.com/SatorNetwork/sator-api/test/framework/user"
	"github.com/SatorNetwork/sator-api/test/framework/utils"
)

func TestCorrectAnswers(t *testing.T) {
	defer app_config.RunAndWait()()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()
	challengeID, err := c.DB.ChallengeDB().ChallengeIDByTitle(context.Background(), "custom2")
	require.NoError(t, err)

	playersNum := 30
	playersInRoom := 10
	users := make([]*user.User, playersNum)
	for i := 0; i < playersNum; i++ {
		users[i] = user.NewInitializedUser(auth.RandomSignUpRequest(), t)
	}

	userExpectedMessages := make([][]*message.Message, playersNum)
	{
		mc := message_container.New(defaultUserExpectedMessages)
		for i := 0; i < playersNum; i++ {
			userExpectedMessages[i] = mc.Copy().Messages()
		}
	}

	messageVerifiers := make([]*message_verifier.MessageVerifier, playersNum)
	for i := 0; i < playersNum; i++ {
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(users[i].AccessToken(), challengeID.String())
		require.NoError(t, err)

		sendMessageSubj := getQuizLinkResp.Data.SendMessageSubj
		recvMessageSubj := getQuizLinkResp.Data.RecvMessageSubj
		userID := getQuizLinkResp.Data.UserID
		serverPublicKey, err := internal_rsa.BytesToPublicKey([]byte(getQuizLinkResp.Data.ServerPublicKey))
		require.NoError(t, err)

		natsSubscriber, err := nats_subscriber.New(userID, sendMessageSubj, recvMessageSubj, t)
		require.NoError(t, err)
		natsSubscriber.SetQuestionMessageCallback(nats_subscriber.ReplyWithCorrectAnswerCallback)
		natsSubscriber.SetEncryptor(envelope.NewEncryptor(serverPublicKey))
		natsSubscriber.SetDecryptor(envelope.NewDecryptor(users[i].PrivateKey()))
		if i == 0 {
			natsSubscriber.EnableDebugMode()
		}
		err = natsSubscriber.Start()
		require.NoError(t, err)
		defer func() {
			err := natsSubscriber.Close()
			require.NoError(t, err)
		}()

		messageVerifiers[i] = message_verifier.New(userExpectedMessages[i], natsSubscriber.GetMessageChan(), t)
		go messageVerifiers[i].Start()
		defer messageVerifiers[i].Close()

		// Wait until room will be closed
		if (i+1)%playersInRoom == 0 {
			time.Sleep(3*time.Second + app_config.AppConfigForTests.QuizLobbyLatency)
		}

		// TODO(evg): investigate and remove
		time.Sleep(time.Second * 2)
	}

	time.Sleep(time.Second * 40)

	for idx, mv := range messageVerifiers {
		fmt.Printf("Idx: %v\n", idx)
		err := mv.NonStrictVerify()
		require.NoError(t, err)
	}
}
