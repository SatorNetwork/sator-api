package quiz_v2

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
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/consts"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

func TestQuizV2Sandbox(t *testing.T) {
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

	signUpRequest2 := auth.RandomSignUpRequest()
	signUpResp2, err := c.Auth.SignUp(signUpRequest2)
	require.NoError(t, err)
	require.NotNil(t, signUpResp2)
	require.NotEmpty(t, signUpResp2.AccessToken)

	err = c.Auth.VerifyAcount(signUpResp2.AccessToken, &auth.VerifyAccountRequest{
		OTP: "12345",
	})
	require.NoError(t, err)

	var (
		user1ExpectedMessages = []*message.Message{
			{
				MessageType: message.PlayerIsJoinedMessageType,
				PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{
					Username: signUpRequest.Username,
				},
			},
		}

		user2ExpectedMessages = []*message.Message{
			{
				MessageType: message.PlayerIsJoinedMessageType,
				PlayerIsJoinedMessage: &message.PlayerIsJoinedMessage{
					Username: signUpRequest2.Username,
				},
			},
			{
				MessageType: message.CountdownMessageType,
				CountdownMessage: &message.CountdownMessage{
					SecondsLeft: 3,
				},
			},
			{
				MessageType: message.CountdownMessageType,
				CountdownMessage: &message.CountdownMessage{
					SecondsLeft: 2,
				},
			},
			{
				MessageType: message.CountdownMessageType,
				CountdownMessage: &message.CountdownMessage{
					SecondsLeft: 1,
				},
			},
			{
				MessageType: message.QuestionMessageType,
				QuestionMessage: &message.QuestionMessage{
					Text: "question1",
				},
			},
			{
				MessageType: message.AnswerReplyMessageType,
				AnswerReplyMessage: &message.AnswerReplyMessage{
					Success: true,
				},
			},
			{
				MessageType: message.QuestionMessageType,
				QuestionMessage: &message.QuestionMessage{
					Text: "question2",
				},
			},
			{
				MessageType: message.AnswerReplyMessageType,
				AnswerReplyMessage: &message.AnswerReplyMessage{
					Success: true,
				},
			},
			{
				MessageType: message.QuestionMessageType,
				QuestionMessage: &message.QuestionMessage{
					Text: "question3",
				},
			},
			{
				MessageType: message.AnswerReplyMessageType,
				AnswerReplyMessage: &message.AnswerReplyMessage{
					Success: true,
				},
			},
		}
	)

	var user1MessageVerifier *message_verifier.MessageVerifier
	{
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(signUpResp.AccessToken, consts.DefaultChallengeID)
		require.NoError(t, err)

		sendMessageSubj := getQuizLinkResp.Data.SendMessageSubj
		recvMessageSubj := getQuizLinkResp.Data.RecvMessageSubj
		userID := getQuizLinkResp.Data.UserID

		natsSubscriber, err := nats_subscriber.New(userID, sendMessageSubj, recvMessageSubj, t)
		require.NoError(t, err)
		natsSubscriber.SetQuestionMessageCallback(nats_subscriber.ReplyWithTrueCallback)
		err = natsSubscriber.Start()
		require.NoError(t, err)
		defer func() {
			err := natsSubscriber.Close()
			require.NoError(t, err)
		}()

		messageVerifier := message_verifier.New(user1ExpectedMessages, natsSubscriber.GetMessageChan())
		go messageVerifier.Start()
		defer messageVerifier.Close()

		time.Sleep(time.Second * 10)

		err = messageVerifier.Verify()
		require.NoError(t, err)

		user1MessageVerifier = messageVerifier
	}

	{
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(signUpResp2.AccessToken, consts.DefaultChallengeID)
		require.NoError(t, err)

		sendMessageSubj := getQuizLinkResp.Data.SendMessageSubj
		recvMessageSubj := getQuizLinkResp.Data.RecvMessageSubj
		userID := getQuizLinkResp.Data.UserID

		natsSubscriber, err := nats_subscriber.New(userID, sendMessageSubj, recvMessageSubj, t)
		require.NoError(t, err)
		natsSubscriber.SetQuestionMessageCallback(nats_subscriber.ReplyWithTrueCallback)
		err = natsSubscriber.Start()
		require.NoError(t, err)
		defer func() {
			err := natsSubscriber.Close()
			require.NoError(t, err)
		}()

		messageVerifier := message_verifier.New(user2ExpectedMessages, natsSubscriber.GetMessageChan())
		go messageVerifier.Start()
		defer messageVerifier.Close()

		time.Sleep(time.Second * 10)

		err = messageVerifier.Verify()
		require.NoError(t, err)
	}

	{
		user1ExpectedMessages = append(user1ExpectedMessages, user2ExpectedMessages...)
		user1MessageVerifier.SetExpectedMessages(user1ExpectedMessages)

		err = user1MessageVerifier.Verify()
		require.NoError(t, err)
	}
}
