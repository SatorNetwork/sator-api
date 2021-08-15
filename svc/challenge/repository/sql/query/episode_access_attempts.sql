-- name: GetAskedQuestionsByEpisodeID :many
SELECT question_id
FROM attempts
WHERE user_id = $1 AND episode_id = $2;
-- name: GetEpisodeIDByQuestionID :one
SELECT episode_id
FROM attempts
WHERE user_id = $1 AND question_id = $2
    LIMIT 1;
-- name: AddAttempt :one
INSERT INTO attempts (user_id, episode_id, question_id, answer_id, valid)
VALUES ($1, $2, $3, $4, $5) RETURNING *;
-- name: CountAttempts :one
SELECT COUNT (*)
FROM attempts
WHERE user_id = $1 AND episode_id = $2 AND created_at > $3 AND answer_id IS NOT NULL;
-- name: UpdateAttempt :exec
UPDATE attempts
SET answer_id = @answer_id, valid = @valid
WHERE user_id = @user_id AND question_id = @question_id;
