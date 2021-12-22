package nats_subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/SatorNetwork/sator-api/internal/encryption/envelope"
	"github.com/pkg/errors"
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
	debugMode               bool
	keepaliveCfg            *KeepaliveCfg

	encryptor *envelope.Encryptor
	decryptor *envelope.Decryptor

	done chan struct{}

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

// TODO: move to constructor
func (s *natsSubscriber) SetEncryptor(encryptor *envelope.Encryptor) {
	s.encryptor = encryptor
}

func (s *natsSubscriber) SetDecryptor(decryptor *envelope.Decryptor) {
	s.decryptor = decryptor
}

func (s *natsSubscriber) Start() error {
	subscription, err := s.nc.Subscribe(s.recvMessageSubj, func(m *nats.Msg) {
		msg, err := s.decodeAndDecrypt(m.Data)
		require.NoError(s.t, err)

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

func (s *natsSubscriber) encodeAndEncrypt(msg *message.Message) ([]byte, error) {
	messageData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	envelope, err := s.encryptor.Encrypt(messageData)
	if err != nil {
		return nil, errors.Wrapf(err, "can't encrypt message with server's public key")
	}
	envelopeData, err := json.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	return envelopeData, nil
}

func (s *natsSubscriber) decodeAndDecrypt(envelopeData []byte) (*message.Message, error) {
	var envelope envelope.Envelope
	if err := json.Unmarshal(envelopeData, &envelope); err != nil {
		return nil, err
	}

	messageData, err := s.decryptor.Decrypt(&envelope)
	if err != nil {
		return nil, errors.Wrapf(err, "can't decrypt message with player's private key")
	}

	if s.debugMode {
		fmt.Printf("Received a message: %s\n", string(messageData))
	}

	msg := new(message.Message)
	if err := msg.UnmarshalJSON(messageData); err != nil {
		return nil, err
	}

	//if s.debugMode {
	//	fmt.Printf("Received a message: %+v\n", msg)
	//}

	return msg, nil
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
	data, err := s.encodeAndEncrypt(msg)
	if err != nil {
		return errors.Wrapf(err, "can't encode & encrypt message")
	}

	if err := s.nc.Publish(s.sendMessageSubj, data); err != nil {
		return err
	}

	return nil
}

func (s *natsSubscriber) GetMessageChan() <-chan *message.Message {
	return s.recvMessageChan
}
