-- name: AddNewPlayer :exec
INSERT INTO quiz_players (quiz_id, user_id, username, status)
VALUES (@quiz_id, @user_id, @username, @status) ON CONFLICT (quiz_id, user_id) DO NOTHING;
-- name: UpdatePlayerStatus :exec
UPDATE quiz_players
SET status = @status
WHERE quiz_id = @quiz_id
    AND user_id = @user_id;
-- name: CountPlayersInQuiz :one
SELECT COUNT(user_id) AS players
FROM quiz_players
WHERE quiz_id = @quiz_id
GROUP BY quiz_id
LIMIT 1;