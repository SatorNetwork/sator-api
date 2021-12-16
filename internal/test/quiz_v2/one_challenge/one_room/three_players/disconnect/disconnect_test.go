package two_players

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/internal/test/framework/client"
	"github.com/SatorNetwork/sator-api/internal/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/internal/test/framework/message_verifier"
	"github.com/SatorNetwork/sator-api/internal/test/framework/nats_subscriber"
	"github.com/SatorNetwork/sator-api/internal/test/framework/utils"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

func TestUserIsDisconnected(t *testing.T) {
	err := utils.BootstrapIfNeeded(context.Background(), t)
	require.NoError(t, err)

	c := client.NewClient()
	challengeID, err := c.DB.ChallengeDB().ChallengeIDByTitle(context.Background(), "custom1")
	require.NoError(t, err)

	playersNum := 3
	signUpRequests := make([]*auth.SignUpRequest, playersNum)
	signUpResponses := make([]*auth.SignUpResponse, playersNum)
	for i := 0; i < playersNum; i++ {
		signUpRequests[i] = auth.RandomSignUpRequest()
		signUpResponses[i], err = c.Auth.SignUp(signUpRequests[i])
		require.NoError(t, err)
		require.NotNil(t, signUpResponses[i])
		require.NotEmpty(t, signUpResponses[i].AccessToken)

		err = c.Auth.VerifyAcount(signUpResponses[i].AccessToken, &auth.VerifyAccountRequest{
			OTP: "12345",
		})
		require.NoError(t, err)
	}

	var user1MessageVerifier *message_verifier.MessageVerifier
	{
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(signUpResponses[0].AccessToken, challengeID.String())
		require.NoError(t, err)

		sendMessageSubj := getQuizLinkResp.Data.SendMessageSubj
		recvMessageSubj := getQuizLinkResp.Data.RecvMessageSubj
		userID := getQuizLinkResp.Data.UserID

		natsSubscriber, err := nats_subscriber.New(userID, sendMessageSubj, recvMessageSubj, t)
		require.NoError(t, err)
		natsSubscriber.SetQuestionMessageCallback(nats_subscriber.ReplyWithCorrectAnswerCallback)
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
					Username: signUpRequests[0].Username,
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
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(signUpResponses[1].AccessToken, challengeID.String())
		require.NoError(t, err)

		sendMessageSubj := getQuizLinkResp.Data.SendMessageSubj
		recvMessageSubj := getQuizLinkResp.Data.RecvMessageSubj
		userID := getQuizLinkResp.Data.UserID

		natsSubscriber, err := nats_subscriber.New(userID, sendMessageSubj, recvMessageSubj, t)
		require.NoError(t, err)
		natsSubscriber.SetQuestionMessageCallback(nats_subscriber.ReplyWithCorrectAnswerCallback)
		natsSubscriber.SetKeepaliveCfg(&nats_subscriber.KeepaliveCfg{
			Disabled: true,
		})
		err = natsSubscriber.Start()
		require.NoError(t, err)
		defer func() {
			err := natsSubscriber.Close()
			require.NoError(t, err)
		}()

		user2ExpectedMessages := []*message.Message{
			{
				MessageType: message.PlayerIsJoinedMessageType,
				PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{
					Username: signUpRequests[1].Username,
				},
			},
		}
		messageVerifier := message_verifier.New(user2ExpectedMessages, natsSubscriber.GetMessageChan(), t)
		go messageVerifier.Start()
		defer messageVerifier.Close()

		time.Sleep(time.Second * 10)

		err = messageVerifier.Verify()
		require.NoError(t, err)
	}

	{
		user1ExpectedMessages := []*message.Message{
			{
				MessageType: message.PlayerIsJoinedMessageType,
				PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{
					Username: signUpRequests[0].Username,
				},
			},
			{
				MessageType: message.PlayerIsJoinedMessageType,
				PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{
					Username: signUpRequests[1].Username,
				},
			},
			{
				MessageType: message.PlayerIsDisconnectedMessageType,
				PlayerIsDisconnectedMessage: &message.PlayerIsDisconnectedMessage{
					Username: signUpRequests[1].Username,
				},
			},
		}
		user1MessageVerifier.SetExpectedMessages(user1ExpectedMessages)

		err := user1MessageVerifier.Verify()
		require.NoError(t, err)
	}

	{
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(signUpResponses[2].AccessToken, challengeID.String())
		require.NoError(t, err)
		_ = getQuizLinkResp
	}
}
