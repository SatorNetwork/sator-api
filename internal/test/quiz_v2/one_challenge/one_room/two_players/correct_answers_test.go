package two_players

import (
	"context"
	"crypto/rsa"
	"github.com/SatorNetwork/sator-api/internal/encryption/envelope"
	internal_rsa "github.com/SatorNetwork/sator-api/internal/encryption/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/internal/test/framework/client"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/internal/test/framework/message_verifier"
	"github.com/SatorNetwork/sator-api/internal/test/framework/nats_subscriber"
	"github.com/SatorNetwork/sator-api/internal/test/framework/utils"
	"github.com/SatorNetwork/sator-api/internal/test/quiz_v2/message_container"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

func TestCorrectAnswers(t *testing.T) {
	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()
	defaultChallengeID, err := c.DB.ChallengeDB().DefaultChallengeID(context.Background())
	require.NoError(t, err)

	signUpRequest := auth.RandomSignUpRequest()
	signUpResp, err := c.Auth.SignUp(signUpRequest)
	require.NoError(t, err)
	require.NotNil(t, signUpResp)
	require.NotEmpty(t, signUpResp.AccessToken)

	err = c.Auth.VerifyAcount(signUpResp.AccessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(t, err)

	var privateKey1 *rsa.PrivateKey
	{
		privateKey, publicKey, err := internal_rsa.GenerateKeyPair(4096)
		require.NoError(t, err)
		publcKeyBytes, err := internal_rsa.PublicKeyToBytes(publicKey)
		require.NoError(t, err)
		err = c.Auth.RegisterPublicKey(signUpResp.AccessToken, &auth.RegisterPublicKeyRequest{
			PublicKey: string(publcKeyBytes),
		})
		require.NoError(t, err)

		privateKey1 = privateKey
	}

	signUpRequest2 := auth.RandomSignUpRequest()
	signUpResp2, err := c.Auth.SignUp(signUpRequest2)
	require.NoError(t, err)
	require.NotNil(t, signUpResp2)
	require.NotEmpty(t, signUpResp2.AccessToken)

	err = c.Auth.VerifyAcount(signUpResp2.AccessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(t, err)

	var privateKey2 *rsa.PrivateKey
	{
		privateKey, publicKey, err := internal_rsa.GenerateKeyPair(4096)
		require.NoError(t, err)
		publcKeyBytes, err := internal_rsa.PublicKeyToBytes(publicKey)
		require.NoError(t, err)
		err = c.Auth.RegisterPublicKey(signUpResp2.AccessToken, &auth.RegisterPublicKeyRequest{
			PublicKey: string(publcKeyBytes),
		})
		require.NoError(t, err)

		privateKey2 = privateKey
	}

	userExpectedMessages := message_container.New(defaultUserExpectedMessages).
		Modify(
			message_container.PFuncIndex(0),
			func(msg *message.Message) {
				msg.PlayerIsJoinedMessage.Username = signUpRequest.Username
			}).
		Modify(
			message_container.PFuncIndex(1),
			func(msg *message.Message) {
				msg.PlayerIsJoinedMessage.Username = signUpRequest2.Username
			}).
		Messages()

	var user1MessageVerifier *message_verifier.MessageVerifier
	{
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(signUpResp.AccessToken, defaultChallengeID.String())
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
		natsSubscriber.SetDecryptor(envelope.NewDecryptor(privateKey1))
		natsSubscriber.EnableDebugMode()
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
					Username: signUpRequest.Username,
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
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(signUpResp2.AccessToken, defaultChallengeID.String())
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
		natsSubscriber.SetDecryptor(envelope.NewDecryptor(privateKey2))
		err = natsSubscriber.Start()
		require.NoError(t, err)
		defer func() {
			err := natsSubscriber.Close()
			require.NoError(t, err)
		}()

		messageVerifier := message_verifier.New(userExpectedMessages, natsSubscriber.GetMessageChan(), t)
		go messageVerifier.Start()
		defer messageVerifier.Close()

		time.Sleep(time.Second * 15)

		err = messageVerifier.Verify()
		require.NoError(t, err)
	}

	{
		user1MessageVerifier.SetExpectedMessages(userExpectedMessages)

		err = user1MessageVerifier.Verify()
		require.NoError(t, err)
	}
}
