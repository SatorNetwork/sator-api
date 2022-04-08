package two_players

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/lib/encryption/envelope"
	internal_rsa "github.com/SatorNetwork/sator-api/lib/encryption/rsa"
	solana_client "github.com/SatorNetwork/sator-api/lib/solana/client"
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
	err = c.DB.QuizV2DB().Repository().CleanUp(context.Background())
	require.NoError(t, err)
	totalRewards := float64(250)

	sc := solana_client.New(app_config.AppConfigForTests.SolanaApiBaseUrl, solana_client.Config{
		SystemProgram:  app_config.AppConfigForTests.SolanaSystemProgram,
		SysvarRent:     app_config.AppConfigForTests.SolanaSysvarRent,
		SysvarClock:    app_config.AppConfigForTests.SolanaSysvarClock,
		SplToken:       app_config.AppConfigForTests.SolanaSplToken,
		StakeProgramID: app_config.AppConfigForTests.SolanaStakeProgramID,
	})

	user1 := user.NewInitializedUser(auth.RandomSignUpRequest(), t)
	user2 := user.NewInitializedUser(auth.RandomSignUpRequest(), t)

	{
		challenge, err := c.ChallengesClient.GetChallengeById(user1.AccessToken(), defaultChallengeID.String())
		require.NoError(t, err)
		require.Equal(t, 2, challenge.Players)
	}
	{
		challenge, err := c.QuizV2Client.GetChallengeById(user1.AccessToken(), defaultChallengeID.String())
		require.NoError(t, err)
		require.Equal(t, 2, challenge.Players)
		require.Equal(t, 0, challenge.RegisteredPlayersInDB)
	}

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

		{
			challenge, err := c.QuizV2Client.GetChallengeById(user1.AccessToken(), defaultChallengeID.String())
			require.NoError(t, err)
			require.Equal(t, 2, challenge.Players)
			require.Equal(t, 1, challenge.RegisteredPlayersInDB)

			require.Equal(t, "250 SAO", challenge.CurrentPrizePool)
			require.Equal(t, float64(250), challenge.CurrentPrizePoolAmount)
		}
		{
			challengesWithPlayer, err := c.QuizV2Client.GetChallengesSortedByPlayers(user1.AccessToken())
			require.NoError(t, err)
			_ = challengesWithPlayer
			// TODO(evg): debug && uncomment
			//require.GreaterOrEqual(t, len(challengesWithPlayer), 1)
			//require.Equal(t, 1, challengesWithPlayer[0].PlayersNumber)
		}
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
		natsSubscriber.SetQuestionMessageCallback(nats_subscriber.ReplyWithCorrectAnswerCallback)
		natsSubscriber.SetEncryptor(envelope.NewEncryptor(serverPublicKey))
		natsSubscriber.SetDecryptor(envelope.NewDecryptor(user2.PrivateKey()))
		err = natsSubscriber.Start()
		require.NoError(t, err)
		defer func() {
			err := natsSubscriber.Close()
			require.NoError(t, err)
		}()

		messageVerifier := message_verifier.New(userExpectedMessages, natsSubscriber.GetMessageChan(), t)
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

	var user1RewardsAmount float64
	{
		rewardsWallet, err := c.Wallet.GetWalletByType(user1.AccessToken(), "rewards")
		require.NoError(t, err)
		rewardsWalletDetails, err := c.Wallet.GetWalletByID(user1.AccessToken(), rewardsWallet.GetDetailsUrl)
		require.NoError(t, err)
		unclaimedCurrency, err := rewardsWalletDetails.FindUnclaimedCurrency()
		require.NoError(t, err)

		user1RewardsAmount = unclaimedCurrency.Amount
	}

	var user2RewardsAmount float64
	{
		rewardsWallet, err := c.Wallet.GetWalletByType(user2.AccessToken(), "rewards")
		require.NoError(t, err)
		rewardsWalletDetails, err := c.Wallet.GetWalletByID(user2.AccessToken(), rewardsWallet.GetDetailsUrl)
		require.NoError(t, err)
		unclaimedCurrency, err := rewardsWalletDetails.FindUnclaimedCurrency()
		require.NoError(t, err)

		user2RewardsAmount = unclaimedCurrency.Amount
	}

	require.Equal(t, totalRewards*1.01, user1RewardsAmount+user2RewardsAmount)

	{
		_, err := c.QuizV2Client.GetQuizLink(user1.AccessToken(), defaultChallengeID.String())
		require.Error(t, err)
		require.Contains(t, err.Error(), "reward has been already received for this challenge")

		_, err = c.QuizV2Client.GetQuizLink(user2.AccessToken(), defaultChallengeID.String())
		require.Error(t, err)
		require.Contains(t, err.Error(), "reward has been already received for this challenge")
	}

	var user1TxHash string
	var user2TxHash string
	{
		user1ClaimRewards, err := c.RewardsClient.ClaimRewards(user1.AccessToken())
		user1TxHash = getTxHashFromTransactionURL(user1ClaimRewards.TransactionURL)
		require.NoError(t, err)
		require.NotEqual(t, "", user1TxHash)

		user2ClaimRewards, err := c.RewardsClient.ClaimRewards(user2.AccessToken())
		user2TxHash = getTxHashFromTransactionURL(user2ClaimRewards.TransactionURL)
		require.NoError(t, err)
		require.NotEqual(t, "", user2TxHash)
	}

	{
		resp, err := sc.GetTransaction(context.Background(), user1TxHash)
		require.NoError(t, err)
		require.Nil(t, resp.Meta.Err)

		resp, err = sc.GetTransaction(context.Background(), user2TxHash)
		require.NoError(t, err)
		require.Nil(t, resp.Meta.Err)
	}
}

func getTxHashFromTransactionURL(transactionURL string) string {
	s := strings.Split(transactionURL, "https://explorer.solana.com/tx/")
	if len(s) != 2 {
		return ""
	}
	s = strings.Split(s[1], "?cluster=devnet")
	if len(s) != 2 {
		return ""
	}
	return s[0]
}
