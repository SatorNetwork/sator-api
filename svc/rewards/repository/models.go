// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Reward struct {
	ID              uuid.UUID      `json:"id"`
	UserID          uuid.UUID      `json:"user_id"`
	RelationID      uuid.NullUUID  `json:"relation_id"`
	Amount          float64        `json:"amount"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
	CreatedAt       time.Time      `json:"created_at"`
	TransactionType int32          `json:"transaction_type"`
	RelationType    sql.NullString `json:"relation_type"`
	TxHash          sql.NullString `json:"tx_hash"`
	Status          int32          `json:"status"`
}
