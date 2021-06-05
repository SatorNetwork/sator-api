-- name: GetShows :many
SELECT *
FROM shows
ORDER BY has_new_episode DESC,
    updated_at DESC,
    created_at DESC
LIMIT $1 OFFSET $2;