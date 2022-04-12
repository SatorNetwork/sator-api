package events

import "github.com/google/uuid"

type Event interface{}

type ForgetRoomEvent struct {
	RoomID uuid.UUID
}

func NewForgetRoomEvent(roomID uuid.UUID) *ForgetRoomEvent {
	return &ForgetRoomEvent{
		RoomID: roomID,
	}
}
