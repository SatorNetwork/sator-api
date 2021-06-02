-- name: AddReward :exec
INSERT INTO rewards (user_id, quiz_id, amount)
VALUES (@user_id, @quiz_id, @amount);
-- name: GetUnWithdrawnRewards :many
SELECT * FROM rewards
WHERE user_id = $1 AND withdrawn = $2;
-- name: Withdraw :exec
UPDATE rewards
SET withdrawn = true
WHERE user_id = $1;