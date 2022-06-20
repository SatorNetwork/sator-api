-- name: GetFlagByKey :one
SELECT *
FROM flags
WHERE key = $1;

-- name: CreateFlag :one
INSERT INTO flags (
    key,
    value
) VALUES (
    @key,
    @value
) ON CONFLICT DO NOTHING RETURNING *;
