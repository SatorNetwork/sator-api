-- name: CreateIAPReceipt :one
INSERT INTO iap_receipts (
    transaction_id,
    receipt_data,
    receipt_in_json,
    user_id
)
VALUES (
    @transaction_id,
    @receipt_data,
    @receipt_in_json,
    @user_id
) RETURNING *;

-- name: GetIAPReceiptByTxID :one
SELECT * FROM iap_receipts
WHERE transaction_id = @transaction_id;
