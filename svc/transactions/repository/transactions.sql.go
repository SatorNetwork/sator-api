// Code generated by sqlc. DO NOT EDIT.
// source: transactions.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const getTransactionByHash = `-- name: GetTransactionByHash :one
SELECT sender_wallet_id, recipient_wallet_id, transaction_hash, amount, updated_at, created_at
FROM transactions
WHERE transaction_hash = $1
LIMIT 1
`

func (q *Queries) GetTransactionByHash(ctx context.Context, transactionHash string) (Transaction, error) {
	row := q.queryRow(ctx, q.getTransactionByHashStmt, getTransactionByHash, transactionHash)
	var i Transaction
	err := row.Scan(
		&i.SenderWalletID,
		&i.RecipientWalletID,
		&i.TransactionHash,
		&i.Amount,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const storeTransactions = `-- name: StoreTransactions :one
INSERT INTO transactions (
    sender_wallet_id,
    recipient_wallet_id,
    transaction_hash,
    amount
)
VALUES (
           $1,
           $2,
           $3,
           $4
       ) RETURNING sender_wallet_id, recipient_wallet_id, transaction_hash, amount, updated_at, created_at
`

type StoreTransactionsParams struct {
	SenderWalletID    uuid.UUID `json:"sender_wallet_id"`
	RecipientWalletID uuid.UUID `json:"recipient_wallet_id"`
	TransactionHash   string    `json:"transaction_hash"`
	Amount            float64   `json:"amount"`
}

func (q *Queries) StoreTransactions(ctx context.Context, arg StoreTransactionsParams) (Transaction, error) {
	row := q.queryRow(ctx, q.storeTransactionsStmt, storeTransactions,
		arg.SenderWalletID,
		arg.RecipientWalletID,
		arg.TransactionHash,
		arg.Amount,
	)
	var i Transaction
	err := row.Scan(
		&i.SenderWalletID,
		&i.RecipientWalletID,
		&i.TransactionHash,
		&i.Amount,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
