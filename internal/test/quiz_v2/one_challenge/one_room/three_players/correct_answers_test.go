package three_players

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
	"github.com/SatorNetwork/sator-api/internal/test/quiz_v2/message_container"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

func TestCorrectAnswers(t *testing.T) {
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

	userExpectedMessages := make([][]*message.Message, playersNum)
	{
		mc := message_container.New(defaultUserExpectedMessages)
		for i := 0; i < playersNum; i++ {
			mc.Modify(
				message_container.PFuncIndex(i),
				func(msg *message.Message) {
					msg.PlayerIsJoinedMessage.Username = signUpRequests[i].Username
				})
		}
		userExpectedMessages[0] = mc.Copy().Messages()
		userExpectedMessages[1] = mc.Copy().Messages()
		userExpectedMessages[2] = mc.Copy().Messages()
	}

	messageVerifiers := make([]*message_verifier.MessageVerifier, playersNum)
	for i := 0; i < playersNum; i++ {
		getQuizLinkResp, err := c.QuizV2Client.GetQuizLink(signUpResponses[i].AccessToken, challengeID.String())
		require.NoError(t, err)

		sendMessageSubj := getQuizLinkResp.Data.SendMessageSubj
		recvMessageSubj := getQuizLinkResp.Data.RecvMessageSubj
		userID := getQuizLinkResp.Data.UserID

		natsSubscriber, err := nats_subscriber.New(userID, sendMessageSubj, recvMessageSubj, t)
		require.NoError(t, err)
		natsSubscriber.SetQuestionMessageCallback(nats_subscriber.ReplyWithCorrectAnswerCallback)
		if i == 2 {
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
	}

	time.Sleep(time.Second * 10)

	err = messageVerifiers[0].NonStrictVerify()
	require.NoError(t, err)
	err = messageVerifiers[1].NonStrictVerify()
	require.NoError(t, err)
	err = messageVerifiers[2].NonStrictVerify()
	require.NoError(t, err)
}
