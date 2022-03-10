package players_map

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/interfaces"
	"github.com/SatorNetwork/sator-api/svc/quiz_v2/player"
	quiz_v2_repository "github.com/SatorNetwork/sator-api/svc/quiz_v2/repository"
)

type PlayersMap struct {
	roomID       uuid.UUID
	challengeID  uuid.UUID
	players      map[string]player.Player
	playersMutex *sync.Mutex
	qr           interfaces.QuizV2Repository
}

func New(roomID uuid.UUID, challengeID uuid.UUID, qr interfaces.QuizV2Repository) *PlayersMap {
	return &PlayersMap{
		roomID:       roomID,
		challengeID:  challengeID,
		players:      make(map[string]player.Player),
		playersMutex: &sync.Mutex{},
		qr:           qr,
	}
}

func (pm *PlayersMap) AddPlayer(p player.Player) error {
	pm.playersMutex.Lock()
	pm.players[p.ID()] = p
	pm.playersMutex.Unlock()

	playerID, err := uuid.Parse(p.ID())
	if err != nil {
		return errors.Wrap(err, "can't parse player ID")
	}

	err = pm.qr.RegisterNewPlayer(context.Background(), quiz_v2_repository.RegisterNewPlayerParams{
		ChallengeID: pm.challengeID,
		UserID:      playerID,
	})
	if err != nil {
		return errors.Wrap(err, "can't register new player")
	}

	return nil
}

func (pm *PlayersMap) RemovePlayerByID(id string) error {
	pm.playersMutex.Lock()
	delete(pm.players, id)
	pm.playersMutex.Unlock()

	playerID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = pm.qr.UnregisterPlayer(context.Background(), quiz_v2_repository.UnregisterPlayerParams{
		ChallengeID: pm.challengeID,
		UserID:      playerID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (pm *PlayersMap) UnregisterAllPlayersFromDB() {
	pm.ExecuteCallback(func(p player.Player) {
		playerUID, err := uuid.Parse(p.ID())
		if err != nil {
			log.Println(err)
		}

		err = pm.qr.UnregisterPlayer(context.Background(), quiz_v2_repository.UnregisterPlayerParams{
			ChallengeID: pm.challengeID,
			UserID:      playerUID,
		})
		if err != nil {
			log.Println(err)
		}
	})
}

func (pm *PlayersMap) GetPlayerByID(id string) player.Player {
	pm.playersMutex.Lock()
	defer pm.playersMutex.Unlock()

	return pm.players[id]
}

func (pm *PlayersMap) GetPlayerByIDNoLock(id string) player.Player {
	return pm.players[id]
}

func (pm *PlayersMap) PlayersNum() int {
	pm.playersMutex.Lock()
	defer pm.playersMutex.Unlock()

	return len(pm.players)
}

func (pm *PlayersMap) PlayersNumInDB() (int, error) {
	playersInRoom, err := pm.qr.CountPlayersInRoom(context.Background(), pm.challengeID)
	if err != nil {
		return 0, err
	}

	return int(playersInRoom), nil
}

func (pm *PlayersMap) ExecuteCallback(callback func(p player.Player)) {
	pm.playersMutex.Lock()
	defer pm.playersMutex.Unlock()

	for _, p := range pm.players {
		callback(p)
	}
}
