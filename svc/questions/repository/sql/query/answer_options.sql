-- name: GetAnswerByID :one
SELECT *
FROM answer_options
WHERE id = $1
LIMIT 1;
-- name: GetAnswersByQuestionID :many
SELECT *
FROM answer_options
WHERE question_id = $1;
-- name: AddQuestionOption :one
INSERT INTO answer_options (question_id, answer_option, is_correct)
VALUES ($1, $2, $3) RETURNING *;
-- name: CheckAnswer :one
SELECT is_correct
FROM answer_options
WHERE id = $1
LIMIT 1;