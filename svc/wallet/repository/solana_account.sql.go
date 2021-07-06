// Code generated by sqlc. DO NOT EDIT.
// source: solana_account.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const addSolanaAccount = `-- name: AddSolanaAccount :one
INSERT INTO solana_accounts (
        public_key,
        private_key
    )
VALUES (
        $1,
        $2
    ) ON CONFLICT (public_key) DO NOTHING RETURNING id, public_key, private_key, status, updated_at, created_at
`

type AddSolanaAccountParams struct {
	PublicKey  string `json:"public_key"`
	PrivateKey []byte `json:"private_key"`
}

func (q *Queries) AddSolanaAccount(ctx context.Context, arg AddSolanaAccountParams) (SolanaAccount, error) {
	row := q.queryRow(ctx, q.addSolanaAccountStmt, addSolanaAccount, arg.PublicKey, arg.PrivateKey)
	var i SolanaAccount
	err := row.Scan(
		&i.ID,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSolanaAccountByID = `-- name: GetSolanaAccountByID :one
SELECT id, public_key, private_key, status, updated_at, created_at
FROM solana_accounts
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetSolanaAccountByID(ctx context.Context, id uuid.UUID) (SolanaAccount, error) {
	row := q.queryRow(ctx, q.getSolanaAccountByIDStmt, getSolanaAccountByID, id)
	var i SolanaAccount
	err := row.Scan(
		&i.ID,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSolanaAccountByUserID = `-- name: GetSolanaAccountByUserID :one
SELECT id, public_key, private_key, status, updated_at, created_at
FROM solana_accounts
WHERE id = (
        SELECT solana_account_id
        FROM wallets
        WHERE user_id = $1
        LIMIT 1
    )
LIMIT 1
`

func (q *Queries) GetSolanaAccountByUserID(ctx context.Context, userID uuid.UUID) (SolanaAccount, error) {
	row := q.queryRow(ctx, q.getSolanaAccountByUserIDStmt, getSolanaAccountByUserID, userID)
	var i SolanaAccount
	err := row.Scan(
		&i.ID,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
