package room

import "github.com/SatorNetwork/sator-api/svc/quiz_v2/player"

type Room interface {
	ChallengeID() string
	AddPlayer(p player.Player)
	IsFull() bool
	Start()
}