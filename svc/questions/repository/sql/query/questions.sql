-- name: GetQuestionByID :one
SELECT *
FROM questions
WHERE id = $1
    LIMIT 1;
-- name: GetQuestionsByChallengeID :many
SELECT *
FROM questions
WHERE challenge_id = $1
ORDER BY quiestion_order ASC
    LIMIT 1;
-- name: AddQuestion :one
INSERT INTO questions (challenge_id, question, question_order)
VALUES ($1, $2, $3) RETURNING *;
