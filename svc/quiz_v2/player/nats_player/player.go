package nats_player

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player"
)

const (
	defaultChanBuffSize = 10
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

	nc *nats.Conn
}

func NewNatsPlayer(userID, challengeID, username, sendMessageSubj, recvMessageSubj string) (player.Player, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	return &natsPlayer{
		id:          userID,
		username:    username,
		challengeID: challengeID,

		sendMessageSubj: sendMessageSubj,
		recvMessageSubj: recvMessageSubj,

		recvMessageChan: make(chan *message.Message, defaultChanBuffSize),

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
		}

		p.recvMessageChan <- &msg
	})
	if err != nil {
		return err
	}
	p.recvMessageSubscription = subscription

	return nil
}

// TODO(evg): close players when room is closed
func (p *natsPlayer) Close() error {
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
