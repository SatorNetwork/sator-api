// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID                 uuid.UUID    `json:"id"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	ActionUrl          string       `json:"action_url"`
	StartsAt           time.Time    `json:"starts_at"`
	EndsAt             time.Time    `json:"ends_at"`
	UpdatedAt          sql.NullTime `json:"updated_at"`
	CreatedAt          time.Time    `json:"created_at"`
	Type               string       `json:"type"`
	TypeSpecificParams string       `json:"type_specific_params"`
}

type ReadAnnouncement struct {
	AnnouncementID uuid.UUID    `json:"announcement_id"`
	UserID         uuid.UUID    `json:"user_id"`
	UpdatedAt      sql.NullTime `json:"updated_at"`
	CreatedAt      time.Time    `json:"created_at"`
}
