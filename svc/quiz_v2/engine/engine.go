package engine

import (
	"log"
	"sync"
	"time"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/restriction_manager"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/room/default_room"
)

type Engine struct {
	newPlayersChan         chan player.Player
	challengeIDToRoom      map[string]room.Room
	challengeIDToRoomMutex *sync.Mutex

	challenges         interfaces.ChallengesService
	stakeLevels        interfaces.StakeLevels
	rewards            interfaces.RewardsService
	qr                 interfaces.QuizV2Repository
	restrictionManager restriction_manager.RestrictionManager

	shuffleQuestions bool
	quizLobbyLatency time.Duration

	done chan struct{}
}

func New(
	challenges interfaces.ChallengesService,
	stakeLevels interfaces.StakeLevels,
	rewards interfaces.RewardsService,
	qr interfaces.QuizV2Repository,
	restrictionManager restriction_manager.RestrictionManager,
	shuffleQuestions bool,
	quizLobbyLatency time.Duration,
) *Engine {
	return &Engine{
		newPlayersChan:         make(chan player.Player),
		challengeIDToRoom:      make(map[string]room.Room, 0),
		challengeIDToRoomMutex: &sync.Mutex{},

		challenges:         challenges,
		stakeLevels:        stakeLevels,
		rewards:            rewards,
		qr:                 qr,
		restrictionManager: restrictionManager,

		shuffleQuestions: shuffleQuestions,
		quizLobbyLatency: quizLobbyLatency,

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

func (e *Engine) GetRoomDetails(challengeID string) (*room.RoomDetails, error) {
	e.challengeIDToRoomMutex.Lock()
	defer e.challengeIDToRoomMutex.Unlock()

	room, ok := e.challengeIDToRoom[challengeID]
	if !ok {
		return nil, NewErrRoomNotFound(challengeID)
	}

	return room.GetRoomDetails()
}

func (e *Engine) getOrCreateRoom(challengeID string) (room.Room, error) {
	e.challengeIDToRoomMutex.Lock()
	defer e.challengeIDToRoomMutex.Unlock()

	if _, ok := e.challengeIDToRoom[challengeID]; !ok {
		room, err := default_room.New(
			challengeID,
			e.challenges,
			e.stakeLevels,
			e.rewards,
			e.qr,
			e.restrictionManager,
			e.shuffleQuestions,
			e.quizLobbyLatency,
		)
		if err != nil {
			return nil, err
		}
		e.challengeIDToRoom[challengeID] = room
		go room.Start()
	}

	return e.challengeIDToRoom[challengeID], nil
}

func (e *Engine) deleteRoom(challengeID string) {
	e.challengeIDToRoomMutex.Lock()
	defer e.challengeIDToRoomMutex.Unlock()

	delete(e.challengeIDToRoom, challengeID)
}
