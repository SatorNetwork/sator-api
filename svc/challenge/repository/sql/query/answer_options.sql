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
WHERE id = $1 AND question_id = $2
LIMIT 1;
-- name: GetAnswersByIDs :many
SELECT *
FROM answer_options
WHERE question_id = ANY(@question_ids::uuid[]);
-- name: DeleteAnswerByID :exec
DELETE FROM answer_options
WHERE id = @id AND question_id = @question_id;
-- name: UpdateAnswer :exec
UPDATE answer_options
SET id = @id,
    question_id = @question_id,
    answer_option = @answer_option,
    is_correct = @is_correct,
    updated_at = @updated_at
WHERE id = @id;