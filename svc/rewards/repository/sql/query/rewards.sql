-- name: AddTransaction :exec
INSERT INTO rewards (
        user_id,
        relation_id,
        relation_type,
        transaction_type,
        amount
    )
VALUES (
        @user_id,
        @relation_id,
        @relation_type,
        @transaction_type,
        @amount
    );

-- name: Withdraw :exec
UPDATE rewards
SET withdrawn = TRUE
WHERE user_id = @user_id
AND transaction_type = 1;

-- name: GetTotalAmount :one
SELECT SUM(amount)::DOUBLE PRECISION
FROM rewards
WHERE user_id = @user_id
AND withdrawn = FALSE
AND transaction_type = 1
GROUP BY user_id;

-- name: GetAmountAvailableToWithdraw :one
SELECT SUM(amount)::DOUBLE PRECISION
FROM rewards
WHERE user_id = @user_id
AND withdrawn = FALSE
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

