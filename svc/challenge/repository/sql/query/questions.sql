-- name: GetQuestionByID :one
SELECT *
FROM questions
WHERE id = $1
LIMIT 1;
-- name: GetQuestionsByChallengeID :many
SELECT *
FROM questions
WHERE challenge_id = $1
ORDER BY question_order ASC;
-- name: AddQuestion :one
INSERT INTO questions (challenge_id, question, question_order)
VALUES ($1, $2, $3) RETURNING *;
-- name: DeleteQuestionByID :exec
DELETE FROM questions
WHERE id = @id;
-- name: UpdateQuestion :exec
UPDATE questions
SET id = @id,
    challenge_id = @challenge_id,
    question = @question,
    question_order = @question_order,
    updated_at = @updated_at
WHERE id = @id;
-- name: GetQuestionsByChallengeIDWithExceptions :many
SELECT *
FROM questions
WHERE challenge_id = @challenge_id AND id != ANY(@question_ids::uuid[]);