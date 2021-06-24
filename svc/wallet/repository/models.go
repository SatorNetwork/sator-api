// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SolanaAccount struct {
	ID          uuid.UUID     `json:"id"`
	AccountType string        `json:"account_type"`
	PublicKey   string        `json:"public_key"`
	PrivateKey  []byte        `json:"private_key"`
	Status      sql.NullInt32 `json:"status"`
	UpdatedAt   sql.NullTime  `json:"updated_at"`
	CreatedAt   time.Time     `json:"created_at"`
}

type Wallet struct {
	ID              uuid.UUID     `json:"id"`
	UserID          uuid.UUID     `json:"user_id"`
	SolanaAccountID uuid.UUID     `json:"solana_account_id"`
	Status          sql.NullInt32 `json:"status"`
	UpdatedAt       sql.NullTime  `json:"updated_at"`
	CreatedAt       time.Time     `json:"created_at"`
	WalletType      string        `json:"wallet_type"`
}
