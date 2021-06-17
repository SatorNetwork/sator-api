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
WHERE user_id = $1
    AND transaction_type = 1;
-- name: GetTotalAmount :one
SELECT SUM(amount)::DOUBLE PRECISION
FROM rewards
WHERE user_id = @user_id
    AND withdrawn = FALSE
    AND transaction_type = 1
GROUP BY user_id;