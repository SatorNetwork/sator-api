-- name: StoreTransactions :one
INSERT INTO transactions (
    sender_wallet_id,
    recipient_wallet_id,
    transaction_hash,
    amount
)
VALUES (
           @sender_wallet_id,
           @recipient_wallet_id,
           @transaction_hash,
           @amount
       ) RETURNING *;
-- name: GetTransactionByHash :one
SELECT *
FROM transactions
WHERE transaction_hash = $1
LIMIT 1;