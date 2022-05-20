// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: rewards.sql

package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const addTransaction = `-- name: AddTransaction :exec
INSERT INTO rewards (
        user_id,
        relation_id,
        relation_type,
        transaction_type,
        amount,
        tx_hash,
        status
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
    )
`

type AddTransactionParams struct {
	UserID          uuid.UUID      `json:"user_id"`
	RelationID      uuid.NullUUID  `json:"relation_id"`
	RelationType    sql.NullString `json:"relation_type"`
	TransactionType int32          `json:"transaction_type"`
	Amount          float64        `json:"amount"`
	TxHash          sql.NullString `json:"tx_hash"`
	Status          string         `json:"status"`
}

func (q *Queries) AddTransaction(ctx context.Context, arg AddTransactionParams) error {
	_, err := q.exec(ctx, q.addTransactionStmt, addTransaction,
		arg.UserID,
		arg.RelationID,
		arg.RelationType,
		arg.TransactionType,
		arg.Amount,
		arg.TxHash,
		arg.Status,
	)
	return err
}

const getAmountAvailableToWithdraw = `-- name: GetAmountAvailableToWithdraw :one
SELECT SUM(amount)::DOUBLE PRECISION
FROM rewards
WHERE user_id = $1
AND status = 'TransactionStatusAvailable'
AND transaction_type = 1
AND created_at < $2
GROUP BY user_id
`

type GetAmountAvailableToWithdrawParams struct {
	UserID       uuid.UUID `json:"user_id"`
	NotAfterDate time.Time `json:"not_after_date"`
}

func (q *Queries) GetAmountAvailableToWithdraw(ctx context.Context, arg GetAmountAvailableToWithdrawParams) (float64, error) {
	row := q.queryRow(ctx, q.getAmountAvailableToWithdrawStmt, getAmountAvailableToWithdraw, arg.UserID, arg.NotAfterDate)
	var column_1 float64
	err := row.Scan(&column_1)
	return column_1, err
}

const getFailedTransactions = `-- name: GetFailedTransactions :many
SELECT id, user_id, relation_id, amount, updated_at, created_at, transaction_type, relation_type, tx_hash, status
FROM rewards
WHERE status = 'TransactionStatusFailed'
AND transaction_type = 1
`

func (q *Queries) GetFailedTransactions(ctx context.Context) ([]Reward, error) {
	rows, err := q.query(ctx, q.getFailedTransactionsStmt, getFailedTransactions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reward
	for rows.Next() {
		var i Reward
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RelationID,
			&i.Amount,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.TransactionType,
			&i.RelationType,
			&i.TxHash,
			&i.Status,
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

const getRequestedTransactions = `-- name: GetRequestedTransactions :many
SELECT id, user_id, relation_id, amount, updated_at, created_at, transaction_type, relation_type, tx_hash, status
FROM rewards
WHERE status = 'TransactionStatusRequested' AND transaction_type = 2 AND created_at <= NOW() - INTERVAL '1 minute'
`

func (q *Queries) GetRequestedTransactions(ctx context.Context) ([]Reward, error) {
	rows, err := q.query(ctx, q.getRequestedTransactionsStmt, getRequestedTransactions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reward
	for rows.Next() {
		var i Reward
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RelationID,
			&i.Amount,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.TransactionType,
			&i.RelationType,
			&i.TxHash,
			&i.Status,
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

const getScannedQRCodeByUserID = `-- name: GetScannedQRCodeByUserID :one
SELECT id, user_id, relation_id, amount, updated_at, created_at, transaction_type, relation_type, tx_hash, status
FROM rewards
WHERE user_id = $1 AND relation_id = $2 AND relation_type =$3
    LIMIT 1
`

type GetScannedQRCodeByUserIDParams struct {
	UserID       uuid.UUID      `json:"user_id"`
	RelationID   uuid.NullUUID  `json:"relation_id"`
	RelationType sql.NullString `json:"relation_type"`
}

func (q *Queries) GetScannedQRCodeByUserID(ctx context.Context, arg GetScannedQRCodeByUserIDParams) (Reward, error) {
	row := q.queryRow(ctx, q.getScannedQRCodeByUserIDStmt, getScannedQRCodeByUserID, arg.UserID, arg.RelationID, arg.RelationType)
	var i Reward
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.RelationID,
		&i.Amount,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.TransactionType,
		&i.RelationType,
		&i.TxHash,
		&i.Status,
	)
	return i, err
}

const getTotalAmount = `-- name: GetTotalAmount :one
SELECT SUM(amount)::DOUBLE PRECISION
FROM rewards
WHERE user_id = $1
AND status = 'TransactionStatusAvailable'
AND transaction_type = 1
GROUP BY user_id
`

func (q *Queries) GetTotalAmount(ctx context.Context, userID uuid.UUID) (float64, error) {
	row := q.queryRow(ctx, q.getTotalAmountStmt, getTotalAmount, userID)
	var column_1 float64
	err := row.Scan(&column_1)
	return column_1, err
}

const getTransactionsByUserIDPaginated = `-- name: GetTransactionsByUserIDPaginated :many
SELECT id, user_id, relation_id, amount, updated_at, created_at, transaction_type, relation_type, tx_hash, status
FROM rewards
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetTransactionsByUserIDPaginatedParams struct {
	UserID uuid.UUID `json:"user_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) GetTransactionsByUserIDPaginated(ctx context.Context, arg GetTransactionsByUserIDPaginatedParams) ([]Reward, error) {
	rows, err := q.query(ctx, q.getTransactionsByUserIDPaginatedStmt, getTransactionsByUserIDPaginated, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Reward
	for rows.Next() {
		var i Reward
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.RelationID,
			&i.Amount,
			&i.UpdatedAt,
			&i.CreatedAt,
			&i.TransactionType,
			&i.RelationType,
			&i.TxHash,
			&i.Status,
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

const requestTransactionsByUserID = `-- name: RequestTransactionsByUserID :exec
UPDATE rewards
SET status = 'TransactionStatusRequested'
WHERE user_id = $1
AND transaction_type = 1
AND status = 'TransactionStatusAvailable'
`

func (q *Queries) RequestTransactionsByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := q.exec(ctx, q.requestTransactionsByUserIDStmt, requestTransactionsByUserID, userID)
	return err
}

const updateTransactionStatusByTxHash = `-- name: UpdateTransactionStatusByTxHash :exec
UPDATE rewards
SET status = $1
WHERE tx_hash = $2
`

type UpdateTransactionStatusByTxHashParams struct {
	Status string         `json:"status"`
	TxHash sql.NullString `json:"tx_hash"`
}

func (q *Queries) UpdateTransactionStatusByTxHash(ctx context.Context, arg UpdateTransactionStatusByTxHashParams) error {
	_, err := q.exec(ctx, q.updateTransactionStatusByTxHashStmt, updateTransactionStatusByTxHash, arg.Status, arg.TxHash)
	return err
}

const updateTransactionTxHash = `-- name: UpdateTransactionTxHash :exec
UPDATE rewards
SET tx_hash = $1, created_at = NOW()
WHERE tx_hash = $2
`

type UpdateTransactionTxHashParams struct {
	TxHashNew sql.NullString `json:"tx_hash_new"`
	TxHash    sql.NullString `json:"tx_hash"`
}

func (q *Queries) UpdateTransactionTxHash(ctx context.Context, arg UpdateTransactionTxHashParams) error {
	_, err := q.exec(ctx, q.updateTransactionTxHashStmt, updateTransactionTxHash, arg.TxHashNew, arg.TxHash)
	return err
}

const withdraw = `-- name: Withdraw :exec
UPDATE rewards
SET status = 'TransactionStatusWithdrawn'
WHERE user_id = $1
AND transaction_type = 1
AND status = 'TransactionStatusRequested'
`

func (q *Queries) Withdraw(ctx context.Context, userID uuid.UUID) error {
	_, err := q.exec(ctx, q.withdrawStmt, withdraw, userID)
	return err
}
