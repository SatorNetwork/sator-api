-- name: GetShows :many
SELECT *
FROM shows
LIMIT $1 OFFSET $2;