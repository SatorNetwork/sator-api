// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RoomPlayer struct {
	ChallengeID uuid.UUID    `json:"challenge_id"`
	UserID      uuid.UUID    `json:"user_id"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	CreatedAt   time.Time    `json:"created_at"`
}