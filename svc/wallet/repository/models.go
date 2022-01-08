// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type EthereumAccount struct {
	ID         uuid.UUID    `json:"id"`
	PublicKey  []byte       `json:"public_key"`
	PrivateKey []byte       `json:"private_key"`
	Address    string       `json:"address"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
	CreatedAt  time.Time    `json:"created_at"`
}

type SolanaAccount struct {
	ID          uuid.UUID     `json:"id"`
	AccountType string        `json:"account_type"`
	PublicKey   string        `json:"public_key"`
	PrivateKey  []byte        `json:"private_key"`
	Status      sql.NullInt32 `json:"status"`
	UpdatedAt   sql.NullTime  `json:"updated_at"`
	CreatedAt   time.Time     `json:"created_at"`
}

type StakeLevel struct {
	ID             uuid.UUID     `json:"id"`
	MinStakeAmount sql.NullInt32 `json:"min_stake_amount"`
	Title          string        `json:"title"`
	Subtitle       string        `json:"subtitle"`
	Multiplier     sql.NullInt32 `json:"multiplier"`
}

type Wallet struct {
	ID                uuid.UUID     `json:"id"`
	UserID            uuid.UUID     `json:"user_id"`
	SolanaAccountID   uuid.UUID     `json:"solana_account_id"`
	Status            sql.NullInt32 `json:"status"`
	UpdatedAt         sql.NullTime  `json:"updated_at"`
	CreatedAt         time.Time     `json:"created_at"`
	WalletType        string        `json:"wallet_type"`
	Sort              int32         `json:"sort"`
	EthereumAccountID uuid.UUID     `json:"ethereum_account_id"`
}
