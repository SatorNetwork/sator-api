// Code generated by sqlc. DO NOT EDIT.
// source: wallet.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const createWallet = `-- name: CreateWallet :one
INSERT INTO wallets (user_id, solana_account_id, ethereum_account_id, wallet_type, sort)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) RETURNING id, user_id, solana_account_id, status, updated_at, created_at, wallet_type, sort, ethereum_account_id
`

type CreateWalletParams struct {
	UserID            uuid.UUID     `json:"user_id"`
	SolanaAccountID   uuid.UUID     `json:"solana_account_id"`
	EthereumAccountID uuid.NullUUID `json:"ethereum_account_id"`
	WalletType        string        `json:"wallet_type"`
	Sort              int32         `json:"sort"`
}

func (q *Queries) CreateWallet(ctx context.Context, arg CreateWalletParams) (Wallet, error) {
	row := q.queryRow(ctx, q.createWalletStmt, createWallet,
		arg.UserID,
		arg.SolanaAccountID,
		arg.EthereumAccountID,
		arg.WalletType,
		arg.Sort,
	)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SolanaAccountID,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.WalletType,
		&i.Sort,
		&i.EthereumAccountID,
	)
	return i, err
}

const deleteWalletByID = `-- name: DeleteWalletByID :exec
DELETE FROM wallets
WHERE id = $1
`

func (q *Queries) DeleteWalletByID(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteWalletByIDStmt, deleteWalletByID, id)
	return err
}

const getWalletByEthereumAccountID = `-- name: GetWalletByEthereumAccountID :one
SELECT id, user_id, solana_account_id, status, updated_at, created_at, wallet_type, sort, ethereum_account_id
FROM wallets
WHERE ethereum_account_id = $1
    LIMIT 1
`

func (q *Queries) GetWalletByEthereumAccountID(ctx context.Context, ethereumAccountID uuid.NullUUID) (Wallet, error) {
	row := q.queryRow(ctx, q.getWalletByEthereumAccountIDStmt, getWalletByEthereumAccountID, ethereumAccountID)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SolanaAccountID,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.WalletType,
		&i.Sort,
		&i.EthereumAccountID,
	)
	return i, err
}

const getWalletByID = `-- name: GetWalletByID :one
SELECT id, user_id, solana_account_id, status, updated_at, created_at, wallet_type, sort, ethereum_account_id
FROM wallets
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetWalletByID(ctx context.Context, id uuid.UUID) (Wallet, error) {
	row := q.queryRow(ctx, q.getWalletByIDStmt, getWalletByID, id)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SolanaAccountID,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.WalletType,
		&i.Sort,
		&i.EthereumAccountID,
	)
	return i, err
}

const getWalletBySolanaAccountID = `-- name: GetWalletBySolanaAccountID :one
SELECT id, user_id, solana_account_id, status, updated_at, created_at, wallet_type, sort, ethereum_account_id
FROM wallets
WHERE solana_account_id = $1
LIMIT 1
`

func (q *Queries) GetWalletBySolanaAccountID(ctx context.Context, solanaAccountID uuid.UUID) (Wallet, error) {
	row := q.queryRow(ctx, q.getWalletBySolanaAccountIDStmt, getWalletBySolanaAccountID, solanaAccountID)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SolanaAccountID,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.WalletType,
		&i.Sort,
		&i.EthereumAccountID,
	)
	return i, err
}

const getWalletByUserIDAndType = `-- name: GetWalletByUserIDAndType :one
SELECT id, user_id, solana_account_id, status, updated_at, created_at, wallet_type, sort, ethereum_account_id
FROM wallets
WHERE user_id = $1 AND wallet_type = $2
ORDER BY sort ASC 
LIMIT 1
`

type GetWalletByUserIDAndTypeParams struct {
	UserID     uuid.UUID `json:"user_id"`
	WalletType string    `json:"wallet_type"`
}

func (q *Queries) GetWalletByUserIDAndType(ctx context.Context, arg GetWalletByUserIDAndTypeParams) (Wallet, error) {
	row := q.queryRow(ctx, q.getWalletByUserIDAndTypeStmt, getWalletByUserIDAndType, arg.UserID, arg.WalletType)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SolanaAccountID,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.WalletType,
		&i.Sort,
		&i.EthereumAccountID,
	)
	return i, err
}

const getWalletsByUserID = `-- name: GetWalletsByUserID :many
SELECT id, user_id, solana_account_id, status, updated_at, created_at, wallet_type, sort, ethereum_account_id
FROM wallets
WHERE user_id = $1
ORDER BY sort ASC
`

func (q *Queries) GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]Wallet, error) {
	rows, err := q.query(ctx, q.getWalletsByUserIDStmt, getWalletsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wallet
	for rows.Next() {
		var i Wallet
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SolanaAccountID,
			&i.Status,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.WalletType,
			&i.Sort,
			&i.EthereumAccountID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
