-- name: AddTransaction :exec
INSERT INTO rewards (user_id, quiz_id, amount, transaction_type)
VALUES (@user_id, @quiz_id, @amount, @transaction_type);
-- name: GetUserRewardsByStatus :many
SELECT *
FROM rewards
WHERE user_id = $1
    AND withdrawn = $2;
-- name: Withdraw :exec
UPDATE rewards
SET withdrawn = TRUE
WHERE user_id = $1;
-- name: GetTotalAmount :one
SELECT SUM(amount)::DOUBLE PRECISION
FROM rewards
WHERE user_id = $1
    AND withdrawn = FALSE
GROUP BY user_id;