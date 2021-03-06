package message_verifier

import (
	"sort"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

type MessageVerifier struct {
	expectedMessages      []*message.Message
	recvMessageChan       <-chan *message.Message
	receivedMessages      []*message.Message
	receivedMessagesMutex *sync.Mutex

	t    *testing.T
	done chan struct{}
}

func New(expectedMessages []*message.Message, recvMessageChan <-chan *message.Message, t *testing.T) *MessageVerifier {
	return &MessageVerifier{
		expectedMessages:      expectedMessages,
		recvMessageChan:       recvMessageChan,
		receivedMessages:      make([]*message.Message, 0),
		receivedMessagesMutex: &sync.Mutex{},

		t:    t,
		done: make(chan struct{}),
	}
}

func (v *MessageVerifier) SetExpectedMessages(expectedMessages []*message.Message) {
	v.expectedMessages = expectedMessages
}

func (v *MessageVerifier) Start() {
LOOP:
	for {
		select {
		case msg := <-v.recvMessageChan:
			//fmt.Printf("MESSAGE: %v\n", msg)
			v.receivedMessagesMutex.Lock()
			v.receivedMessages = append(v.receivedMessages, msg)
			v.receivedMessagesMutex.Unlock()

		case <-v.done:
			break LOOP
		}
	}
}

func (v *MessageVerifier) Close() {
	v.done <- struct{}{}
}

func (v *MessageVerifier) Verify() error {
	v.receivedMessagesMutex.Lock()
	defer v.receivedMessagesMutex.Unlock()

	if len(v.expectedMessages) != len(v.receivedMessages) {
		return errors.Errorf("expected %v messages, got: %v", len(v.expectedMessages), len(v.receivedMessages))
	}

	messagesNum := len(v.expectedMessages)
	for i := 0; i < messagesNum; i++ {
		expectedMsg := v.expectedMessages[i]
		receivedMsg := v.receivedMessages[i]

		v.compareMessages(expectedMsg, receivedMsg)
	}

	return nil
}

func (v *MessageVerifier) NonStrictVerify() error {
	v.receivedMessagesMutex.Lock()
	defer v.receivedMessagesMutex.Unlock()

	//if len(v.expectedMessages) != len(v.receivedMessages) {
	//	return errors.Errorf("expected %v messages, got: %v", len(v.expectedMessages), len(v.receivedMessages))
	//}

	emap := make(map[message.MessageType]int, 0)
	rmap := make(map[message.MessageType]int, 0)
	for _, emsg := range v.expectedMessages {
		if emsg.MessageType == message.PlayerIsJoinedMessageType {
			continue
		}
		emap[emsg.MessageType]++
	}
	for _, rmsg := range v.receivedMessages {
		if rmsg.MessageType == message.PlayerIsJoinedMessageType {
			continue
		}
		rmap[rmsg.MessageType]++
	}
	require.Equal(v.t, emap, rmap)

	// TODO(evg): enhance non-strict-verification?

	return nil
}

func (v *MessageVerifier) compareMessages(emsg, rmsg *message.Message) {
	require.Equal(v.t, emsg.MessageType, rmsg.MessageType)

	messageType := emsg.MessageType
	switch messageType {
	case message.PlayerIsJoinedMessageType:
		v.comparePlayerIsJoinedMessages(emsg, rmsg)
	case message.CountdownMessageType:
		v.compareCountdownMessages(emsg, rmsg)
	case message.QuestionMessageType:
		v.compareQuestionMessages(emsg, rmsg)
	case message.AnswerMessageType:
	case message.AnswerReplyMessageType:
		v.compareAnswerReplyMessages(emsg, rmsg)
	case message.WinnersTableMessageType:
		v.compareWinnersTableMessages(emsg, rmsg)
	case message.PlayerIsActiveMessageType:
	case message.PlayerIsDisconnectedMessageType:
		v.comparePlayerIsDisconnectedMessages(emsg, rmsg)
	case message.TimeOutMessageType:
		v.compareTimeOutMessages(emsg, rmsg)

	default:
		v.t.Fatalf("<unknown message type>")
	}
}

func (v *MessageVerifier) comparePlayerIsJoinedMessages(emsg, rmsg *message.Message) {
	require.NotNil(v.t, emsg.PlayerIsJoinedMessage)
	require.NotNil(v.t, rmsg.PlayerIsJoinedMessage)
	require.Equal(v.t, emsg.PlayerIsJoinedMessage.Username, rmsg.PlayerIsJoinedMessage.Username)
}

func (v *MessageVerifier) compareCountdownMessages(emsg, rmsg *message.Message) {
	require.NotNil(v.t, emsg.CountdownMessage)
	require.NotNil(v.t, rmsg.CountdownMessage)
	require.Equal(v.t, emsg.CountdownMessage.SecondsLeft, rmsg.CountdownMessage.SecondsLeft)
}

func (v *MessageVerifier) compareQuestionMessages(emsg, rmsg *message.Message) {
	require.NotNil(v.t, emsg.QuestionMessage)
	require.NotNil(v.t, rmsg.QuestionMessage)

	require.Equal(v.t, emsg.QuestionMessage.QuestionText, rmsg.QuestionMessage.QuestionText)
	require.Equal(v.t, emsg.QuestionMessage.TimeForAnswer, rmsg.QuestionMessage.TimeForAnswer)
	require.Equal(v.t, emsg.QuestionMessage.QuestionNumber, rmsg.QuestionMessage.QuestionNumber)
	require.Equal(v.t, emsg.QuestionMessage.TotalQuestions, rmsg.QuestionMessage.TotalQuestions)
	require.Equal(v.t, len(emsg.QuestionMessage.AnswerOptions), len(rmsg.QuestionMessage.AnswerOptions))

	eOptions := emsg.QuestionMessage.AnswerOptions
	rOptions := rmsg.QuestionMessage.AnswerOptions
	sort.Slice(eOptions, func(i, j int) bool {
		return eOptions[i].AnswerText < eOptions[j].AnswerText
	})
	sort.Slice(rOptions, func(i, j int) bool {
		return rOptions[i].AnswerText < rOptions[j].AnswerText
	})
	optionsNum := len(eOptions)
	for i := 0; i < optionsNum; i++ {
		require.Equal(v.t, eOptions[i].AnswerText, rOptions[i].AnswerText)
	}
}

func (v *MessageVerifier) compareAnswerReplyMessages(emsg, rmsg *message.Message) {
	require.NotNil(v.t, emsg.AnswerReplyMessage)
	require.NotNil(v.t, rmsg.AnswerReplyMessage)
	require.Equal(v.t, emsg.AnswerReplyMessage.Success, rmsg.AnswerReplyMessage.Success)
	require.Equal(v.t, emsg.AnswerReplyMessage.SegmentNum, rmsg.AnswerReplyMessage.SegmentNum)
}

func (v *MessageVerifier) compareWinnersTableMessages(emsg, rmsg *message.Message) {
	// TODO(evg): high to predict who will get extra points for fastest answer due to concurrency
	// (but it affects prize pool distribution) so skipping this checking for now
	// require.Equal(v.t, emsg.WinnersTableMessage.PrizePoolDistribution, rmsg.WinnersTableMessage.PrizePoolDistribution)
}

func (v *MessageVerifier) comparePlayerIsDisconnectedMessages(emsg, rmsg *message.Message) {
	require.NotNil(v.t, emsg.PlayerIsDisconnectedMessage)
	require.NotNil(v.t, rmsg.PlayerIsDisconnectedMessage)
	require.Equal(v.t, emsg.PlayerIsDisconnectedMessage.Username, rmsg.PlayerIsDisconnectedMessage.Username)
}

func (v *MessageVerifier) compareTimeOutMessages(emsg, rmsg *message.Message) {
	require.NotNil(v.t, emsg.TimeOutMessage)
	require.NotNil(v.t, rmsg.TimeOutMessage)
	require.Equal(v.t, emsg.TimeOutMessage.Message, rmsg.TimeOutMessage.Message)
}
