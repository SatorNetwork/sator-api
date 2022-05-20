-- name: AddTransaction :exec
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
        @user_id,
        @relation_id,
        @relation_type,
        @transaction_type,
        @amount,
        @tx_hash,
        @status
    );

-- name: Withdraw :exec
UPDATE rewards
SET status = 'TransactionStatusWithdrawn'
WHERE user_id = @user_id
AND transaction_type = 1
AND status = 'TransactionStatusRequested';

-- name: GetTotalAmount :one
SELECT SUM(amount)::DOUBLE PRECISION
FROM rewards
WHERE user_id = @user_id
AND status = 'TransactionStatusAvailable'
AND transaction_type = 1
GROUP BY user_id;

-- name: GetAmountAvailableToWithdraw :one
SELECT SUM(amount)::DOUBLE PRECISION
FROM rewards
WHERE user_id = @user_id
AND status = 'TransactionStatusAvailable'
AND transaction_type = 1
AND created_at < @not_after_date
GROUP BY user_id;

-- name: GetTransactionsByUserIDPaginated :many
SELECT *
FROM rewards
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetScannedQRCodeByUserID :one
SELECT *
FROM rewards
WHERE user_id = $1 AND relation_id = $2 AND relation_type =$3
    LIMIT 1;

-- name: RequestTransactionsByUserID :exec
UPDATE rewards
SET status = 'TransactionStatusRequested'
WHERE user_id = @user_id
AND transaction_type = 1
AND status = 'TransactionStatusAvailable';

-- name: UpdateTransactionStatusByTxHash :exec
UPDATE rewards
SET status = @status
WHERE tx_hash = @tx_hash;

-- name: GetFailedTransactions :many
SELECT *
FROM rewards
WHERE status = 'TransactionStatusFailed'
AND transaction_type = 1;

-- name: GetRequestedTransactions :many
SELECT *
FROM rewards
WHERE status = 'TransactionStatusRequested' AND transaction_type = 2 AND created_at <= NOW() - INTERVAL '1 minute';

-- name: UpdateTransactionTxHash :exec
UPDATE rewards
SET tx_hash = @tx_hash_new, created_at = NOW()
WHERE tx_hash = @tx_hash;
