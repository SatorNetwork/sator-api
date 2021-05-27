-- name: GetQuestionByID :one
SELECT *
FROM questions
WHERE id = $1
    LIMIT 1;
-- name: GetQuestionByChallengeID :many
SELECT *
FROM questions
WHERE challenge_id = $1
    LIMIT 1;
-- name: AddQuestion :one
INSERT INTO questions (id, challenge_id, question, question_order)
VALUES ($1, $2, $3, $4) RETURNING *;
