// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: iap_receipts.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const createIAPReceipt = `-- name: CreateIAPReceipt :one
INSERT INTO iap_receipts (
    transaction_id,
    receipt_data,
    receipt_in_json,
    user_id
)
VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING id, transaction_id, receipt_data, receipt_in_json, user_id, updated_at, created_at
`

type CreateIAPReceiptParams struct {
	TransactionID string    `json:"transaction_id"`
	ReceiptData   string    `json:"receipt_data"`
	ReceiptInJson string    `json:"receipt_in_json"`
	UserID        uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateIAPReceipt(ctx context.Context, arg CreateIAPReceiptParams) (IapReceipt, error) {
	row := q.queryRow(ctx, q.createIAPReceiptStmt, createIAPReceipt,
		arg.TransactionID,
		arg.ReceiptData,
		arg.ReceiptInJson,
		arg.UserID,
	)
	var i IapReceipt
	err := row.Scan(
		&i.ID,
		&i.TransactionID,
		&i.ReceiptData,
		&i.ReceiptInJson,
		&i.UserID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getIAPReceiptByTxID = `-- name: GetIAPReceiptByTxID :one
SELECT id, transaction_id, receipt_data, receipt_in_json, user_id, updated_at, created_at FROM iap_receipts
WHERE transaction_id = $1
`

func (q *Queries) GetIAPReceiptByTxID(ctx context.Context, transactionID string) (IapReceipt, error) {
	row := q.queryRow(ctx, q.getIAPReceiptByTxIDStmt, getIAPReceiptByTxID, transactionID)
	var i IapReceipt
	err := row.Scan(
		&i.ID,
		&i.TransactionID,
		&i.ReceiptData,
		&i.ReceiptInJson,
		&i.UserID,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
