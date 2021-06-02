-- name: AddReward :exec
INSERT INTO rewards (user_id, quiz_id, amount)
VALUES (@user_id, @quiz_id, @amount);