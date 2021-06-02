package quiz

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/SatorNetwork/sator-api/svc/questions"
	"github.com/dustin/go-broadcast"
	"github.com/google/uuid"
)

type (
	// Quiz hub
	Hub struct {
		QuizID      uuid.UUID
		ChallengeID uuid.UUID

		challenge challenge.Challenge

		sendMsg            broadcast.Broadcaster
		sendQuestionResult broadcast.Broadcaster

		questions       map[string]questions.Question
		questionsSentAt map[string]time.Time

		players map[string]*PlayerHub
	}

	PlayerHub struct {
		UserID   uuid.UUID
		Username string

		sendMsg     broadcast.Broadcaster
		receiveAnsw broadcast.Broadcaster
	}
)

// Setup new player hub
func NewPlayerHub(userID uuid.UUID, username string) *PlayerHub {
	return &PlayerHub{
		UserID:      userID,
		Username:    username,
		sendMsg:     broadcast.NewBroadcaster(10),
		receiveAnsw: broadcast.NewBroadcaster(10),
	}
}

// Close player hub
func (ph *PlayerHub) Close() error {
	if err := ph.sendMsg.Close(); err != nil {
		return fmt.Errorf("could not close sending message broadcast: %w", err)
	}

	if err := ph.receiveAnsw.Close(); err != nil {
		return fmt.Errorf("could not close received answers broadcast: %w", err)
	}

	return nil
}

// Setup new quiz hub
func NewQuizHub() *Hub {
	return &Hub{}
}

// Adds player to a quiz hub annd send other players user_cconnected message
func (h *Hub) AddPlayer(userID uuid.UUID, username string) error {
	if _, ok := h.players[userID.String()]; !ok {
		h.players[userID.String()] = NewPlayerHub(userID, username)

		if err := h.SendMessage(UserConnectedMessage, User{
			UserID:   userID.String(),
			Username: username,
		}); err != nil {
			return fmt.Errorf("add player: could not encode message: %w", err)
		}
	}

	return nil
}

func (h *Hub) RemovePlayer(userID uuid.UUID) error {
	if p, ok := h.players[userID.String()]; !ok {
		if err := p.Close(); err != nil {
			return fmt.Errorf("remove player: could not close player hub: %w", err)
		}

		if err := h.SendMessage(UserDisonnectedMessage, User{
			UserID:   userID.String(),
			Username: p.Username,
		}); err != nil {
			return fmt.Errorf("remove player: could not encode message: %w", err)
		}
	}

	return nil
}

// Sends message to general quiz channel
func (h *Hub) SendMessage(msgType string, msg interface{}) error {
	b, err := json.Marshal(Message{
		Type:    msgType,
		SentAt:  time.Now(),
		Payload: msg,
	})
	if err != nil {
		return fmt.Errorf("could not encode message: %w", err)
	}
	h.sendMsg.Submit(b)
	return nil
}

// Sends message to general quiz channel
func (h *Hub) SendPersonalMessage(userID uuid.UUID, msgType string, msg interface{}) error {
	b, err := json.Marshal(Message{
		Type:    msgType,
		SentAt:  time.Now(),
		Payload: msg,
	})
	if err != nil {
		return fmt.Errorf("could not encode message: %w", err)
	}

	if ph, ok := h.players[userID.String()]; ok {
		ph.sendMsg.Submit(b)
		return nil
	}

	return fmt.Errorf("player with id=%s not found", userID.String())
}
