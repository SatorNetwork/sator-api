-- name: GetAnswerByID :one
SELECT *
FROM question_options
WHERE id = $1
    LIMIT 1;
-- name: GetAnswersByQuestionID :many
SELECT *
FROM question_options
WHERE question_id = $1
    LIMIT 1;
-- name: AddQuestionOption :one
INSERT INTO question_options (question_id, question_option, is_correct)
VALUES ($1, $2, $3) RETURNING *;
-- name: CheckAnswer :one
SELECT is_correct
FROM question_options
WHERE id = $1 AND question_id = $2
    LIMIT 1;