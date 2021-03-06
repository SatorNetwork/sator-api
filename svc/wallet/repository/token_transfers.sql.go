// Code generated by sqlc. DO NOT EDIT.
// source: token_transfers.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addTokenTransfer = `-- name: AddTokenTransfer :one
INSERT INTO token_transfers (user_id, sender_address, recipient_address, amount, tx_hash, status)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6
    ) RETURNING id, user_id, sender_address, recipient_address, tx_hash, amount, status, updated_at, created_at
`

type AddTokenTransferParams struct {
	UserID           uuid.UUID      `json:"user_id"`
	SenderAddress    string         `json:"sender_address"`
	RecipientAddress string         `json:"recipient_address"`
	Amount           float64        `json:"amount"`
	TxHash           sql.NullString `json:"tx_hash"`
	Status           int32          `json:"status"`
}

func (q *Queries) AddTokenTransfer(ctx context.Context, arg AddTokenTransferParams) (TokenTransfer, error) {
	row := q.queryRow(ctx, q.addTokenTransferStmt, addTokenTransfer,
		arg.UserID,
		arg.SenderAddress,
		arg.RecipientAddress,
		arg.Amount,
		arg.TxHash,
		arg.Status,
	)
	var i TokenTransfer
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SenderAddress,
		&i.RecipientAddress,
		&i.TxHash,
		&i.Amount,
		&i.Status,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const checkRecipientAddress = `-- name: CheckRecipientAddress :one
SELECT count(DISTINCT user_id)
FROM token_transfers 
WHERE recipient_address = $1
    AND user_id != $2
`

type CheckRecipientAddressParams struct {
	RecipientAddress string    `json:"recipient_address"`
	UserID           uuid.UUID `json:"user_id"`
}

func (q *Queries) CheckRecipientAddress(ctx context.Context, arg CheckRecipientAddressParams) (int64, error) {
	row := q.queryRow(ctx, q.checkRecipientAddressStmt, checkRecipientAddress, arg.RecipientAddress, arg.UserID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const doesUserHaveFraudulentTransfers = `-- name: DoesUserHaveFraudulentTransfers :one
SELECT (count(DISTINCT user_id) > 0)::BOOLEAN as fraud_detected
FROM token_transfers 
WHERE user_id = $1
    AND status = 3
`

func (q *Queries) DoesUserHaveFraudulentTransfers(ctx context.Context, userID uuid.UUID) (bool, error) {
	row := q.queryRow(ctx, q.doesUserHaveFraudulentTransfersStmt, doesUserHaveFraudulentTransfers, userID)
	var fraud_detected bool
	err := row.Scan(&fraud_detected)
	return fraud_detected, err
}

const doesUserMakeTransferForLastMinute = `-- name: DoesUserMakeTransferForLastMinute :one
SELECT (count(*) > 0)::BOOLEAN as found_transfer
FROM token_transfers 
WHERE user_id = $1
    AND created_at > now() - interval '1 minute'
`

func (q *Queries) DoesUserMakeTransferForLastMinute(ctx context.Context, userID uuid.UUID) (bool, error) {
	row := q.queryRow(ctx, q.doesUserMakeTransferForLastMinuteStmt, doesUserMakeTransferForLastMinute, userID)
	var found_transfer bool
	err := row.Scan(&found_transfer)
	return found_transfer, err
}

const updateTokenTransfer = `-- name: UpdateTokenTransfer :exec
UPDATE token_transfers
SET status = $1,
    tx_hash = $2
WHERE id = $3
`

type UpdateTokenTransferParams struct {
	Status int32          `json:"status"`
	TxHash sql.NullString `json:"tx_hash"`
	ID     uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateTokenTransfer(ctx context.Context, arg UpdateTokenTransferParams) error {
	_, err := q.exec(ctx, q.updateTokenTransferStmt, updateTokenTransfer, arg.Status, arg.TxHash, arg.ID)
	return err
}
