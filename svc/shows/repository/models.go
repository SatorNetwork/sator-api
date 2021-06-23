// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Show struct {
	ID            uuid.UUID      `json:"id"`
	Title         string         `json:"title"`
	Cover         string         `json:"cover"`
	HasNewEpisode bool           `json:"has_new_episode"`
	UpdatedAt     sql.NullTime   `json:"updated_at"`
	CreatedAt     time.Time      `json:"created_at"`
	Category      sql.NullString `json:"category"`
}
