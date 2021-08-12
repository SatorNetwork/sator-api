-- name: GetAttemptByEpisodeID :one
SELECT *
FROM attempts
WHERE user_id = $1 AND episode_id = $2
    LIMIT 1;
-- name: GetAttemptByQuestionID :one
SELECT *
FROM attempts
WHERE user_id = $1 AND question_id = $2
    LIMIT 1;
-- name: AddAttempt :one
INSERT INTO attempts (user_id, episode_id, question_id, answer_id, valid)
VALUES ($1, $2, $3, $4, $5) RETURNING *;
-- name: DeleteAttempt :exec
DELETE FROM attempts
WHERE user_id = $1 AND episode_id = $2;
-- name: UpdateAttempt :exec
UPDATE attempts
SET question_id = @question_id,
    answer_id = @answer_id,
    valid = @valid
WHERE user_id = @user_id AND episode_id = @episode_id;
-- name: CountAttempts :one
SELECT COUNT (*)
FROM attempts
WHERE user_id = $1 AND episode_id = $2 AND created_at > $3;
