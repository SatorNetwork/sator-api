package nats_player

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/SatorNetwork/sator-api/internal/encryption/envelope"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player/nats_player/status_transactor"
)

const (
	defaultChanBuffSize = 10
	disconnectTimeout   = 3 * time.Second
)

type natsPlayer struct {
	id          string
	username    string
	challengeID string

	sendMessageSubj           string
	outgoingMessagesBuffer    []*message.Message
	outgoingMessagesBufferMtx *sync.Mutex
	recvMessageSubj           string
	recvMessageChan           chan *message.Message
	// It will be set during Start method
	recvMessageSubscription *nats.Subscription

	statusIsUpdatedChan      chan struct{}
	st                       *status_transactor.StatusTransactor
	lastUserIsActiveMsg      time.Time
	lastUserIsActiveMsgMutex *sync.Mutex
	connectChan              chan struct{}
	disconnectChan           chan struct{}

	encryptor *envelope.Encryptor

	done chan struct{}

	nc *nats.Conn
}

func NewNatsPlayer(
	userID string,
	challengeID string,
	username string,
	natsURL string,
	sendMessageSubj string,
	recvMessageSubj string,
	playerPublicKey *rsa.PublicKey,
) (player.Player, error) {
	statusIsUpdatedChan := make(chan struct{}, defaultChanBuffSize)

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, errors.Wrapf(err, "can't connect to nats server, url is %v", natsURL)
	}

	return &natsPlayer{
		id:          userID,
		username:    username,
		challengeID: challengeID,

		sendMessageSubj:           sendMessageSubj,
		outgoingMessagesBuffer:    make([]*message.Message, 0),
		outgoingMessagesBufferMtx: &sync.Mutex{},
		recvMessageSubj:           recvMessageSubj,
		recvMessageChan:           make(chan *message.Message, defaultChanBuffSize),

		statusIsUpdatedChan:      statusIsUpdatedChan,
		st:                       status_transactor.New(statusIsUpdatedChan),
		lastUserIsActiveMsgMutex: &sync.Mutex{},
		connectChan:              make(chan struct{}, defaultChanBuffSize),
		disconnectChan:           make(chan struct{}, defaultChanBuffSize),

		encryptor: envelope.NewEncryptor(playerPublicKey),

		done: make(chan struct{}),

		nc: nc,
	}, nil
}

func (p *natsPlayer) ID() string {
	return p.id
}

func (p *natsPlayer) Username() string {
	return p.username
}

func (p *natsPlayer) ChallengeID() string {
	return p.challengeID
}

func (p *natsPlayer) Start() error {
	subscription, err := p.nc.Subscribe(p.recvMessageSubj, func(m *nats.Msg) {
		var msg message.Message
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			log.Printf("can't unmarshal nats message: %v\n", err)
			return
		}

		if err := msg.CheckConsistency(); err != nil {
			log.Println(err)
			return
		}

		if msg.MessageType == message.PlayerIsActiveMessageType {
			p.lastUserIsActiveMsgMutex.Lock()
			p.lastUserIsActiveMsg = time.Now()
			p.lastUserIsActiveMsgMutex.Unlock()

			p.st.SetStatus(status_transactor.PlayerConnectedStatus)
			return
		}

		p.recvMessageChan <- &msg
	})
	if err != nil {
		return err
	}
	p.recvMessageSubscription = subscription

	go p.startEventProcessor()

	return nil
}

// NOTE: should be run as a goroutine
func (p *natsPlayer) startEventProcessor() {
	ticker := time.NewTicker(disconnectTimeout)
LOOP:
	for {
		select {
		case <-ticker.C:
			if p.checkIfDisconnected() {
				p.st.SetStatus(status_transactor.PlayerDisconnectedStatus)
			}

		case <-p.statusIsUpdatedChan:
			switch p.st.GetStatus() {
			case status_transactor.PlayerDisconnectedStatus:
				p.disconnectChan <- struct{}{}
			case status_transactor.PlayerConnectedStatus:
				p.flushOutgoingBuffer()
				p.connectChan <- struct{}{}
			}

		case <-p.done:
			ticker.Stop()
			break LOOP
		}
	}
}

func (p *natsPlayer) Close() error {
	close(p.done)

	return p.recvMessageSubscription.Unsubscribe()
}

func (p *natsPlayer) SendMessage(msg *message.Message) error {
	if p.st.GetStatus() == status_transactor.UndefinedStatus {
		p.outgoingMessagesBufferMtx.Lock()
		p.outgoingMessagesBuffer = append(p.outgoingMessagesBuffer, msg)
		p.outgoingMessagesBufferMtx.Unlock()
		return nil
	}

	data, err := p.encodeAndEncrypt(msg)
	if err != nil {
		return errors.Wrap(err, "can't encode & encrypt message")
	}
	if err := p.nc.Publish(p.sendMessageSubj, data); err != nil {
		return err
	}

	return nil
}

func (p *natsPlayer) encodeAndEncrypt(msg *message.Message) ([]byte, error) {
	messageData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	envelope, err := p.encryptor.Encrypt(messageData)
	if err != nil {
		return nil, errors.Wrapf(err, "can't encrypt message with player's public key")
	}
	envelopeData, err := json.Marshal(envelope)
	if err != nil {
		return nil, err
	}

	return envelopeData, nil
}

func (p *natsPlayer) flushOutgoingBuffer() {
	p.outgoingMessagesBufferMtx.Lock()
	defer p.outgoingMessagesBufferMtx.Unlock()

	for _, msg := range p.outgoingMessagesBuffer {
		if err := p.SendMessage(msg); err != nil {
			log.Printf("can't send message: %v\n", err)
		}
	}

	p.outgoingMessagesBuffer = make([]*message.Message, 0)
}

func (p *natsPlayer) GetMessageStream() <-chan *message.Message {
	return p.recvMessageChan
}

func (p *natsPlayer) checkIfDisconnected() bool {
	p.lastUserIsActiveMsgMutex.Lock()
	defer p.lastUserIsActiveMsgMutex.Unlock()

	return time.Now().Sub(p.lastUserIsActiveMsg).Nanoseconds() > disconnectTimeout.Nanoseconds()
}

func (p *natsPlayer) ConnectChan() <-chan struct{} {
	return p.connectChan
}

func (p *natsPlayer) DisconnectChan() <-chan struct{} {
	return p.disconnectChan
}
