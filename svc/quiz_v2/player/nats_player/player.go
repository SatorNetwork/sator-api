package nats_player

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/encryption/envelope"
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
	avatar      string
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
	decryptor *envelope.Decryptor

	done chan struct{}

	nc *nats.Conn
}

func NewNatsPlayer(
	userID string,
	challengeID string,
	username string,
	avatar string,
	natsURL string,
	sendMessageSubj string,
	recvMessageSubj string,
	playerPublicKey *rsa.PublicKey,
	serverPrivateKey *rsa.PrivateKey,
) (player.Player, error) {
	statusIsUpdatedChan := make(chan struct{}, defaultChanBuffSize)

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, errors.Wrapf(err, "can't connect to nats server, url is %v", natsURL)
	}

	return &natsPlayer{
		id:          userID,
		username:    username,
		avatar:      avatar,
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
		decryptor: envelope.NewDecryptor(serverPrivateKey),

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

func (p *natsPlayer) Avatar() string {
	return p.avatar
}

func (p *natsPlayer) ChallengeID() string {
	return p.challengeID
}

func (p *natsPlayer) Start() error {
	subscription, err := p.nc.Subscribe(p.recvMessageSubj, func(m *nats.Msg) {
		{
			tmpl := `
				Received encrypted message
				User's ID: %v
				Username:  %v
				Base64-encoded encrypted message: %v
				`
			log.Printf(tmpl, p.id, p.username, base64.StdEncoding.EncodeToString(m.Data))
		}

		msg, err := p.decodeAndDecrypt(m.Data)
		if err != nil {
			log.Printf("can't decode & decrypt message: %v\n", err)
			return
		}
		{
			tmpl := `
				Message successfully decrypted & decoded
				User's ID: %v
				Username:  %v
				Decrypted & decoded message: %+v
				`
			log.Printf(tmpl, p.id, p.username, msg)
		}

		if err := msg.CheckConsistency(); err != nil {
			tmpl := `
				Message isn't consistent: %v
				User's ID: %v
				Username:  %v
				Decrypted & decoded message: %+v
				`
			log.Printf(tmpl, err, p.id, p.username, msg)
			return
		}

		if msg.MessageType == message.PlayerIsActiveMessageType {
			p.lastUserIsActiveMsgMutex.Lock()
			p.lastUserIsActiveMsg = time.Now()
			p.lastUserIsActiveMsgMutex.Unlock()

			tmpl := `
				Player status is connected
				User's ID: %v
				Username:  %v
				`
			log.Printf(tmpl, p.id, p.username)

			p.st.SetStatus(status_transactor.PlayerConnectedStatus)
			return
		}

		p.recvMessageChan <- msg
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
				tmpl := `
					Player status is disconnected
					User's ID: %v
					Username:  %v
					`
				log.Printf(tmpl, p.id, p.username)

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
	tmpl := `
		Going to send message
		User's ID: %v
		Username:  %v
		Message:   %+v
		`
	log.Printf(tmpl, p.id, p.username, msg)

	if p.st.GetStatus() == status_transactor.UndefinedStatus {
		tmpl := `
		Going to send message (player status is Undefined); keep message in outgoingMessagesBuffer
		User's ID: %v
		Username:  %v
		Message:   %+v
		`
		log.Printf(tmpl, p.id, p.username, msg)

		p.outgoingMessagesBufferMtx.Lock()
		p.outgoingMessagesBuffer = append(p.outgoingMessagesBuffer, msg)
		p.outgoingMessagesBufferMtx.Unlock()
		return nil
	}

	data, err := p.encodeAndEncrypt(msg)
	if err != nil {
		err := errors.Wrap(err, "can't encode & encrypt message")
		log.Println(err)
		return err
	}
	if err := p.nc.Publish(p.sendMessageSubj, data); err != nil {
		err := errors.Wrapf(err, "can't publish message to %v", p.sendMessageSubj)
		log.Println(err)
		return err
	}

	return nil
}

func (p *natsPlayer) encodeAndEncrypt(msg *message.Message) ([]byte, error) {
	messageData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	{
		tmpl := `
			Going to send message
			User's ID: %v
			Username:  %v
			Encoded message: %s
			`
		log.Printf(tmpl, p.id, p.username, messageData)
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

func (p *natsPlayer) decodeAndDecrypt(envelopeData []byte) (*message.Message, error) {
	var envelope envelope.Envelope
	if err := json.Unmarshal(envelopeData, &envelope); err != nil {
		return nil, err
	}

	messageData, err := p.decryptor.Decrypt(&envelope)
	if err != nil {
		return nil, errors.Wrapf(err, "can't decrypt message with server's private key")
	}

	{
		tmpl := `
			Message successfully decrypted
			User's ID: %v
			Username:  %v
			Decrypted message: %s
			`
		log.Printf(tmpl, p.id, p.username, messageData)
	}

	msg := new(message.Message)
	if err := msg.UnmarshalJSON(messageData); err != nil {
		return nil, err
	}

	return msg, nil
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
