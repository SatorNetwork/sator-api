package engine

import (
	"fmt"
)

type ErrRoomNotFound struct {
	challengeID string
}

func NewErrRoomNotFound(challengeID string) *ErrRoomNotFound {
	return &ErrRoomNotFound{
		challengeID: challengeID,
	}
}

func (e *ErrRoomNotFound) Error() string {
	return fmt.Sprintf("room for challenge %v doesn't exist", e.challengeID)
}
