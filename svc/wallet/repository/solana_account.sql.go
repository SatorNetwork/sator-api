// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: solana_account.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const addSolanaAccount = `-- name: AddSolanaAccount :one
INSERT INTO solana_accounts (
        account_type,
        public_key,
        private_key
    )
VALUES (
        $1,
        $2,
        $3
    ) ON CONFLICT (public_key) DO NOTHING RETURNING id, account_type, public_key, private_key, status, updated_at, created_at
`

type AddSolanaAccountParams struct {
	AccountType string `json:"account_type"`
	PublicKey   string `json:"public_key"`
	PrivateKey  []byte `json:"private_key"`
}

func (q *Queries) AddSolanaAccount(ctx context.Context, arg AddSolanaAccountParams) (SolanaAccount, error) {
	row := q.queryRow(ctx, q.addSolanaAccountStmt, addSolanaAccount, arg.AccountType, arg.PublicKey, arg.PrivateKey)
	var i SolanaAccount
	err := row.Scan(
		&i.ID,
		&i.AccountType,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSolanaAccountByID = `-- name: GetSolanaAccountByID :one
SELECT id, account_type, public_key, private_key, status, updated_at, created_at
FROM solana_accounts
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetSolanaAccountByID(ctx context.Context, id uuid.UUID) (SolanaAccount, error) {
	row := q.queryRow(ctx, q.getSolanaAccountByIDStmt, getSolanaAccountByID, id)
	var i SolanaAccount
	err := row.Scan(
		&i.ID,
		&i.AccountType,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSolanaAccountByType = `-- name: GetSolanaAccountByType :one
SELECT id, account_type, public_key, private_key, status, updated_at, created_at
FROM solana_accounts
WHERE account_type = $1
LIMIT 1
`

func (q *Queries) GetSolanaAccountByType(ctx context.Context, accountType string) (SolanaAccount, error) {
	row := q.queryRow(ctx, q.getSolanaAccountByTypeStmt, getSolanaAccountByType, accountType)
	var i SolanaAccount
	err := row.Scan(
		&i.ID,
		&i.AccountType,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSolanaAccountByUserIDAndType = `-- name: GetSolanaAccountByUserIDAndType :one
SELECT id, account_type, public_key, private_key, status, updated_at, created_at
FROM solana_accounts
WHERE id = (
        SELECT solana_account_id
        FROM wallets
        WHERE user_id = $1
            AND wallet_type = $2
        LIMIT 1
    )
LIMIT 1
`

type GetSolanaAccountByUserIDAndTypeParams struct {
	UserID     uuid.UUID `json:"user_id"`
	WalletType string    `json:"wallet_type"`
}

func (q *Queries) GetSolanaAccountByUserIDAndType(ctx context.Context, arg GetSolanaAccountByUserIDAndTypeParams) (SolanaAccount, error) {
	row := q.queryRow(ctx, q.getSolanaAccountByUserIDAndTypeStmt, getSolanaAccountByUserIDAndType, arg.UserID, arg.WalletType)
	var i SolanaAccount
	err := row.Scan(
		&i.ID,
		&i.AccountType,
		&i.PublicKey,
		&i.PrivateKey,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getSolanaAccountTypeByPublicKey = `-- name: GetSolanaAccountTypeByPublicKey :one
SELECT account_type
FROM solana_accounts
WHERE public_key = $1
LIMIT 1
`

func (q *Queries) GetSolanaAccountTypeByPublicKey(ctx context.Context, publicKey string) (string, error) {
	row := q.queryRow(ctx, q.getSolanaAccountTypeByPublicKeyStmt, getSolanaAccountTypeByPublicKey, publicKey)
	var account_type string
	err := row.Scan(&account_type)
	return account_type, err
}
