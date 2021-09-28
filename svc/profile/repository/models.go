// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
	CreatedAt time.Time      `json:"created_at"`
	Avatar    sql.NullString `json:"avatar"`
}
