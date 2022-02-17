package room

import "github.com/SatorNetwork/sator-api/svc/quiz_v2/player"

type RoomDetails struct {
	PlayersToStart    int32
	RegisteredPlayers int
}

type Room interface {
	ChallengeID() string
	AddPlayer(p player.Player)
	IsFull() bool
	Start()
	GetRoomDetails() *RoomDetails
}
