// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID             uuid.UUID    `json:"id"`
	Email          string       `json:"email"`
	InvitedBy      uuid.UUID    `json:"invited_by"`
	InvitedAt      time.Time    `json:"invited_at"`
	AcceptedBy     uuid.UUID    `json:"accepted_by"`
	AcceptedAt     sql.NullTime `json:"accepted_at"`
	RewardReceived sql.NullBool `json:"reward_received"`
}
