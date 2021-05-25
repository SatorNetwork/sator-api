-- name: GetShows :many
SELECT *
FROM shows
ORDER BY updated_at,
    created_at DESC
LIMIT $1 OFFSET $2;