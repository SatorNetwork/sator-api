package nats_subscriber

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

const (
	defaultChanBuffSize = 10
)

type messageCallback func(s *natsSubscriber, msg *message.Message)

type natsSubscriber struct {
	nc              *nats.Conn
	userID          string
	sendMessageSubj string
	recvMessageSubj string
	recvMessageChan chan *message.Message
	// It will be set during Start method
	recvMessageSubscription *nats.Subscription

	questionMessageCallback messageCallback

	debugMode bool

	t *testing.T
}

func New(userID, sendMessageSubj, recvMessageSubj string, t *testing.T) (*natsSubscriber, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	return &natsSubscriber{
		nc:              nc,
		userID:          userID,
		sendMessageSubj: sendMessageSubj,
		recvMessageSubj: recvMessageSubj,
		recvMessageChan: make(chan *message.Message, defaultChanBuffSize),
		t:               t,
	}, nil
}

func (s *natsSubscriber) SetQuestionMessageCallback(cb messageCallback) {
	s.questionMessageCallback = cb
}

func (s *natsSubscriber) IsDebugModeEnabled() bool {
	return s.debugMode
}

func (s *natsSubscriber) EnableDebugMode() {
	s.debugMode = true
}

func (s *natsSubscriber) Start() error {
	subscription, err := s.nc.Subscribe(s.recvMessageSubj, func(m *nats.Msg) {
		if s.debugMode {
			fmt.Printf("Received a message: %s\n", string(m.Data))
		}

		var msg message.Message
		err := json.Unmarshal(m.Data, &msg)
		require.NoError(s.t, err)

		s.recvMessageChan <- &msg

		switch msg.MessageType {
		case message.QuestionMessageType:
			if s.questionMessageCallback != nil {
				s.questionMessageCallback(s, &msg)
			}
		}
	})
	if err != nil {
		return err
	}
	s.recvMessageSubscription = subscription

	return nil
}

func (s *natsSubscriber) Close() error {
	return s.recvMessageSubscription.Unsubscribe()
}

func (s *natsSubscriber) SendMessage(msg *message.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := s.nc.Publish(s.sendMessageSubj, data); err != nil {
		return err
	}

	return nil
}

func (s *natsSubscriber) GetMessageChan() <-chan *message.Message {
	return s.recvMessageChan
}
