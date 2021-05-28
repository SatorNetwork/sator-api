// Code generated by sqlc. DO NOT EDIT.
// source: wallet.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createWallet = `-- name: CreateWallet :one
INSERT INTO wallets (
        user_id,
        asset_name,
        wallet_address,
        public_key,
        private_key,
        status
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    ) RETURNING id, user_id, asset_name, wallet_address, public_key, private_key, status, updated_at, created_at
`

type CreateWalletParams struct {
	UserID        uuid.UUID     `json:"user_id"`
	AssetName     string        `json:"asset_name"`
	WalletAddress string        `json:"wallet_address"`
	PublicKey     string        `json:"public_key"`
	PrivateKey    []byte        `json:"private_key"`
	Status        sql.NullInt32 `json:"status"`
}

func (q *Queries) CreateWallet(ctx context.Context, arg CreateWalletParams) (Wallet, error) {
	row := q.queryRow(ctx, q.createWalletStmt, createWallet,
		arg.UserID,
		arg.AssetName,
		arg.WalletAddress,
		arg.PublicKey,
		arg.PrivateKey,
		arg.Status,
	)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AssetName,
		&i.WalletAddress,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getWalletByAssetName = `-- name: GetWalletByAssetName :one
SELECT id, user_id, asset_name, wallet_address, public_key, private_key, status, updated_at, created_at
FROM wallets
WHERE user_id = $1
    AND asset_name = $2
LIMIT 1
`

type GetWalletByAssetNameParams struct {
	UserID    uuid.UUID `json:"user_id"`
	AssetName string    `json:"asset_name"`
}

func (q *Queries) GetWalletByAssetName(ctx context.Context, arg GetWalletByAssetNameParams) (Wallet, error) {
	row := q.queryRow(ctx, q.getWalletByAssetNameStmt, getWalletByAssetName, arg.UserID, arg.AssetName)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.AssetName,
		&i.WalletAddress,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getWalletByID = `-- name: GetWalletByID :one
SELECT id, user_id, asset_name, wallet_address, public_key, private_key, status, updated_at, created_at
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
		&i.AssetName,
		&i.WalletAddress,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getWalletsByUserID = `-- name: GetWalletsByUserID :many
SELECT id, user_id, asset_name, wallet_address, public_key, private_key, status, updated_at, created_at
FROM wallets
WHERE user_id = $1
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
			&i.AssetName,
			&i.WalletAddress,
			&i.PublicKey,
			&i.PrivateKey,
			&i.Status,
			&i.UpdatedAt,
			&i.CreatedAt,
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
