package player

import (
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

//go:generate mockgen -destination=mock_player/player.go -package=mock_player github.com/SatorNetwork/sator-api/svc/quiz_v2/player Player
type Player interface {
	ID() string
	Username() string
	ChallengeID() string
	Start() error
	SendMessage(msg *message.Message) error
	GetMessageStream() <-chan *message.Message
	Close() error
	ConnectionNotifier
}

type ConnectionNotifier interface {
	ConnectChan() <-chan struct{}
	DisconnectChan() <-chan struct{}
}
