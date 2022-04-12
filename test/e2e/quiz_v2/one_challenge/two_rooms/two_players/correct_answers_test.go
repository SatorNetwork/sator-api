package two_players

import (
	"context"
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
	defaultChallengeID, err := c.DB.ChallengeDB().DefaultChallengeID(context.Background())
	require.NoError(t, err)

	playersNum := 4
	playersInRoom := 2
	users := make([]*user.User, playersNum)
	for i := 0; i < playersNum; i++ {
		users[i] = user.NewInitializedUser(auth.RandomSignUpRequest(), t)
	}

	userExpectedMessages := make([][]*message.Message, playersNum)
	{
		mc := message_container.New(defaultUserExpectedMessages).
			Modify(
				message_container.PFuncIndex(0),
				func(msg *message.Message) {
					msg.PlayerIsJoinedMessage.Username = users[0].Username()
				}).
			Modify(
				message_container.PFuncIndex(1),
				func(msg *message.Message) {
					msg.PlayerIsJoinedMessage.Username = users[1].Username()
				})
		userExpectedMessages[0] = mc.Copy().Messages()
		userExpectedMessages[1] = mc.Copy().Messages()
	}

	{
		mc := message_container.New(defaultUserExpectedMessages).
			Modify(
				message_container.PFuncIndex(0),
				func(msg *message.Message) {
					msg.PlayerIsJoinedMessage.Username = users[2].Username()
				}).
			Modify(
				message_container.PFuncIndex(1),
				func(msg *message.Message) {
					msg.PlayerIsJoinedMessage.Username = users[3].Username()
				})
		userExpectedMessages[2] = mc.Copy().Messages()
		userExpectedMessages[3] = mc.Copy().Messages()
	}

	messageVerifiers := make([]*message_verifier.MessageVerifier, playersNum)
	for i := 0; i < playersNum; i++ {
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(users[i].AccessToken(), defaultChallengeID.String())
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

	time.Sleep(time.Second * 25)

	//for _, mv := range messageVerifiers {
	//	err := mv.Verify()
	//	require.NoError(t, err)
	//}

	err = messageVerifiers[0].Verify()
	require.NoError(t, err)
	err = messageVerifiers[1].Verify()
	require.NoError(t, err)
	err = messageVerifiers[2].Verify()
	require.NoError(t, err)
	err = messageVerifiers[3].Verify()
	require.NoError(t, err)
}
