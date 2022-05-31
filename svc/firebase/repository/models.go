// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type DisabledNotification struct {
	UserID    uuid.UUID    `json:"user_id"`
	Topic     string       `json:"topic"`
	UpdatedAt sql.NullTime `json:"updated_at"`
	CreatedAt time.Time    `json:"created_at"`
}

type FirebaseRegistrationToken struct {
	DeviceID          string       `json:"device_id"`
	UserID            uuid.UUID    `json:"user_id"`
	RegistrationToken string       `json:"registration_token"`
	UpdatedAt         sql.NullTime `json:"updated_at"`
	CreatedAt         time.Time    `json:"created_at"`
}
