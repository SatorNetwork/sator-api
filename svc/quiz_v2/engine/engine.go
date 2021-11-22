package engine

import (
	"time"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room"
)

type Engine struct {
	newPlayersChan    chan player.Player
	challengeIDToRoom map[string]room.Room

	done chan struct{}
}

func New() *Engine {
	return &Engine{
		newPlayersChan:    make(chan player.Player),
		challengeIDToRoom: make(map[string]room.Room, 0),

		done: make(chan struct{}),
	}
}

func (e *Engine) Start() {
LOOP:
	for {
		select {
		case newPlayer := <-e.newPlayersChan:
			room := e.getOrCreateRoom(newPlayer.ChallengeID())

			room.AddPlayer(newPlayer)
			// wait for some time until user will be properly registered in room
			// TODO(evg): need smth better; add SyncAddPlayer API?
			time.Sleep(time.Second)
			if room.IsFull() {
				e.deleteRoom(room.ChallengeID())
			}

		case <-e.done:
			break LOOP
		}
	}
}

func (e *Engine) Close() {
	close(e.done)
}

func (e *Engine) AddPlayer(p player.Player) {
	e.newPlayersChan <- p
}

func (e *Engine) getOrCreateRoom(challengeID string) room.Room {
	if _, ok := e.challengeIDToRoom[challengeID]; !ok {
		room := default_room.New(challengeID)
		e.challengeIDToRoom[challengeID] = room
		go room.Start()
	}

	return e.challengeIDToRoom[challengeID]
}

func (e *Engine) deleteRoom(challengeID string) {
	delete(e.challengeIDToRoom, challengeID)
}
