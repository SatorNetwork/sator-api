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

func TestNoAnswers(t *testing.T) {
	defer app_config.RunAndWait()()

	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()
	defaultChallengeID, err := c.DB.ChallengeDB().DefaultChallengeID(context.Background())
	require.NoError(t, err)

	user1 := user.NewInitializedUser(auth.RandomSignUpRequest(), t)
	user2 := user.NewInitializedUser(auth.RandomSignUpRequest(), t)

	userExpectedMessages := message_container.New(defaultUserExpectedMessages).
		Modify(
			message_container.PFuncIndex(0),
			func(msg *message.Message) {
				msg.PlayerIsJoinedMessage.Username = user1.Username()
			}).
		Modify(
			message_container.PFuncIndex(1),
			func(msg *message.Message) {
				msg.PlayerIsJoinedMessage.Username = user2.Username()
			}).
		Messages()

	var user1MessageVerifier *message_verifier.MessageVerifier
	{
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(user1.AccessToken(), defaultChallengeID.String())
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
		natsSubscriber.SetDecryptor(envelope.NewDecryptor(user1.PrivateKey()))
		err = natsSubscriber.Start()
		require.NoError(t, err)
		defer func() {
			err := natsSubscriber.Close()
			require.NoError(t, err)
		}()

		user1ExpectedMessages := []*message.Message{
			{
				MessageType: message.PlayerIsJoinedMessageType,
				PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{
					Username: user1.Username(),
				},
			},
		}
		messageVerifier := message_verifier.New(user1ExpectedMessages, natsSubscriber.GetMessageChan(), t)
		go messageVerifier.Start()
		defer messageVerifier.Close()

		time.Sleep(time.Second * 10)

		err = messageVerifier.Verify()
		require.NoError(t, err)

		user1MessageVerifier = messageVerifier
	}

	{
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(user2.AccessToken(), defaultChallengeID.String())
		require.NoError(t, err)

		sendMessageSubj := getQuizLinkResp.Data.SendMessageSubj
		recvMessageSubj := getQuizLinkResp.Data.RecvMessageSubj
		userID := getQuizLinkResp.Data.UserID
		serverPublicKey, err := internal_rsa.BytesToPublicKey([]byte(getQuizLinkResp.Data.ServerPublicKey))
		require.NoError(t, err)

		natsSubscriber, err := nats_subscriber.New(userID, sendMessageSubj, recvMessageSubj, t)
		require.NoError(t, err)
		natsSubscriber.SetQuestionMessageCallback(nats_subscriber.NoAnswerCallback)
		natsSubscriber.SetEncryptor(envelope.NewEncryptor(serverPublicKey))
		natsSubscriber.SetDecryptor(envelope.NewDecryptor(user2.PrivateKey()))
		err = natsSubscriber.Start()
		require.NoError(t, err)
		defer func() {
			err := natsSubscriber.Close()
			require.NoError(t, err)
		}()

		timeOutMessage, err := message.NewTimeOutMessage(&message.TimeOutMessage{
			Message: "time is over",
		})
		require.NoError(t, err)

		user2ExpectedMessages := message_container.New(userExpectedMessages).
			Replace(
				message_container.PFuncMessageType(message.AnswerReplyMessageType),
				timeOutMessage,
			).
			Messages()

		messageVerifier := message_verifier.New(user2ExpectedMessages, natsSubscriber.GetMessageChan(), t)
		go messageVerifier.Start()
		defer messageVerifier.Close()

		time.Sleep(time.Second*25 + app_config.AppConfigForTests.QuizLobbyLatency)

		err = messageVerifier.Verify()
		require.NoError(t, err)
	}

	{
		user1MessageVerifier.SetExpectedMessages(userExpectedMessages)

		err = user1MessageVerifier.Verify()
		require.NoError(t, err)
	}
}
