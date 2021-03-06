// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Blacklist struct {
	RestrictedType  string `json:"restricted_type"`
	RestrictedValue string `json:"restricted_value"`
}

type User struct {
	ID             uuid.UUID      `json:"id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	Password       []byte         `json:"password"`
	Disabled       bool           `json:"disabled"`
	VerifiedAt     sql.NullTime   `json:"verified_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	CreatedAt      time.Time      `json:"created_at"`
	Role           string         `json:"role"`
	BlockReason    sql.NullString `json:"block_reason"`
	SanitizedEmail sql.NullString `json:"sanitized_email"`
	EmailHash      sql.NullString `json:"email_hash"`
	KycStatus      sql.NullString `json:"kyc_status"`
	PublicKey      sql.NullString `json:"public_key"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
}

type UserVerification struct {
	RequestType      int32     `json:"request_type"`
	UserID           uuid.UUID `json:"user_id"`
	Email            string    `json:"email"`
	VerificationCode []byte    `json:"verification_code"`
	CreatedAt        time.Time `json:"created_at"`
}

type UsersDevice struct {
	UserID   uuid.UUID `json:"user_id"`
	DeviceID string    `json:"device_id"`
}

type Whitelist struct {
	AllowedType  string `json:"allowed_type"`
	AllowedValue string `json:"allowed_value"`
}
