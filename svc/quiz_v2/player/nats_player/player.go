package nats_player

import (
	"encoding/json"
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

	sendMessageSubj string
	recvMessageSubj string
	recvMessageChan chan *message.Message
	// It will be set during Start method
	recvMessageSubscription *nats.Subscription

	statusIsUpdatedChan      chan struct{}
	st                       *status_transactor.StatusTransactor
	lastUserIsActiveMsg      time.Time
	lastUserIsActiveMsgMutex *sync.Mutex
	connectChan              chan struct{}
	disconnectChan           chan struct{}

	done chan struct{}

	nc *nats.Conn
}

func NewNatsPlayer(userID, challengeID, username, natsURL, sendMessageSubj, recvMessageSubj string) (player.Player, error) {
	statusIsUpdatedChan := make(chan struct{}, defaultChanBuffSize)

	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, errors.Wrapf(err, "can't connect to nats server, url is %v", natsURL)
	}

	return &natsPlayer{
		id:          userID,
		username:    username,
		challengeID: challengeID,

		sendMessageSubj: sendMessageSubj,
		recvMessageSubj: recvMessageSubj,

		recvMessageChan: make(chan *message.Message, defaultChanBuffSize),

		statusIsUpdatedChan:      statusIsUpdatedChan,
		st:                       status_transactor.New(statusIsUpdatedChan),
		lastUserIsActiveMsgMutex: &sync.Mutex{},
		connectChan:              make(chan struct{}, defaultChanBuffSize),
		disconnectChan:           make(chan struct{}, defaultChanBuffSize),

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
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := p.nc.Publish(p.sendMessageSubj, data); err != nil {
		return err
	}

	return nil
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
