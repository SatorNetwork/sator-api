package engine

import (
	"log"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/restriction_manager"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room"
)

type Engine struct {
	newPlayersChan    chan player.Player
	ChallengeIDToRoom map[string]room.Room

	challenges         interfaces.ChallengesService
	stakeLevels        interfaces.StakeLevels
	rewards            interfaces.RewardsService
	restrictionManager restriction_manager.RestrictionManager

	shuffleQuestions bool

	done chan struct{}
}

func New(
	challenges interfaces.ChallengesService,
	stakeLevels interfaces.StakeLevels,
	rewards interfaces.RewardsService,
	restrictionManager restriction_manager.RestrictionManager,
	shuffleQuestions bool,
) *Engine {
	return &Engine{
		newPlayersChan:    make(chan player.Player),
		ChallengeIDToRoom: make(map[string]room.Room, 0),

		challenges:         challenges,
		stakeLevels:        stakeLevels,
		rewards:            rewards,
		restrictionManager: restrictionManager,

		shuffleQuestions: shuffleQuestions,

		done: make(chan struct{}),
	}
}

func (e *Engine) Start() {
LOOP:
	for {
		select {
		case newPlayer := <-e.newPlayersChan:
			room, err := e.getOrCreateRoom(newPlayer.ChallengeID())
			if err != nil {
				log.Println(err)
				continue
			}

			room.AddPlayer(newPlayer)
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

func (e *Engine) getOrCreateRoom(challengeID string) (room.Room, error) {
	if _, ok := e.ChallengeIDToRoom[challengeID]; !ok {
		room, err := default_room.New(challengeID, e.challenges, e.stakeLevels, e.rewards, e.restrictionManager, e.shuffleQuestions)
		if err != nil {
			return nil, err
		}
		e.ChallengeIDToRoom[challengeID] = room
		go room.Start()
	}

	return e.ChallengeIDToRoom[challengeID], nil
}

func (e *Engine) deleteRoom(challengeID string) {
	delete(e.ChallengeIDToRoom, challengeID)
}
