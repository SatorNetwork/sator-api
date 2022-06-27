-- name: GetFlagByKey :one
SELECT *
FROM flags
WHERE key = $1;

-- name: CreateFlag :exec
INSERT INTO flags (
    key,
    value
) VALUES (
    @key,
    @value
) ON CONFLICT DO NOTHING;

-- name: UpdateFlag :one
UPDATE flags SET value = @value
WHERE key = @key
RETURNING *;

-- name: GetFlags :many
SELECT * FROM flags;
