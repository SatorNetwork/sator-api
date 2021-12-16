package nats_subscriber

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

const (
	defaultChanBuffSize = 10
)

type messageCallback func(s *natsSubscriber, msg *message.Message)

type KeepaliveCfg struct {
	Disabled bool
}

var defaultKeepaliveCfg = &KeepaliveCfg{
	Disabled: false,
}

type natsSubscriber struct {
	nc              *nats.Conn
	userID          string
	sendMessageSubj string
	recvMessageSubj string
	recvMessageChan chan *message.Message
	// It will be set during Start method
	recvMessageSubscription *nats.Subscription

	questionMessageCallback messageCallback

	debugMode    bool
	keepaliveCfg *KeepaliveCfg
	done         chan struct{}

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
		keepaliveCfg:    defaultKeepaliveCfg,
		done:            make(chan struct{}),
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

func (s *natsSubscriber) SetKeepaliveCfg(keepaliveCfg *KeepaliveCfg) {
	s.keepaliveCfg = keepaliveCfg
}

func (s *natsSubscriber) Start() error {
	subscription, err := s.nc.Subscribe(s.recvMessageSubj, func(m *nats.Msg) {
		if s.debugMode {
			fmt.Printf("Received a message: %s\n", string(m.Data))
		}

		msg := new(message.Message)
		err := msg.UnmarshalJSON(m.Data)
		require.NoError(s.t, err)

		//if s.debugMode {
		//	fmt.Printf("Received a message: %+v\n", msg)
		//}

		s.recvMessageChan <- msg

		switch msg.MessageType {
		case message.QuestionMessageType:
			if s.questionMessageCallback != nil {
				s.questionMessageCallback(s, msg)
			}
		}
	})
	if err != nil {
		return err
	}
	s.recvMessageSubscription = subscription

	go s.startEventProcessor()

	return nil
}

// NOTE: should be run as a goroutine
func (s *natsSubscriber) startEventProcessor() {
	ticker := time.NewTicker(time.Second)
LOOP:
	for {
		select {
		case <-ticker.C:
			if s.keepaliveCfg.Disabled {
				continue
			}

			payload := message.PlayerIsActiveMessage{}
			respMsg, err := message.NewPlayerIsActiveMessage(&payload)
			require.NoError(s.t, err)
			err = s.SendMessage(respMsg)
			require.NoError(s.t, err)

		case <-s.done:
			ticker.Stop()
			break LOOP
		}
	}
}

func (s *natsSubscriber) Close() error {
	close(s.done)

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
