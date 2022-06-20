// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type WatcherTransaction struct {
	ID                     uuid.UUID    `json:"id"`
	SerializedMessage      string       `json:"serialized_message"`
	LatestValidBlockHeight int64        `json:"latest_valid_block_height"`
	AccountAliases         []string     `json:"account_aliases"`
	TxHash                 string       `json:"tx_hash"`
	Status                 string       `json:"status"`
	UpdatedAt              sql.NullTime `json:"updated_at"`
	CreatedAt              time.Time    `json:"created_at"`
	Retries                int32        `json:"retries"`
}