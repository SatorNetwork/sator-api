-- name: AddReward :exec
INSERT INTO rewards (user_id, quiz_id, amount)
VALUES (@user_id, @quiz_id, @amount);
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