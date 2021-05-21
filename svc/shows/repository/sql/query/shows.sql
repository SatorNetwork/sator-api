-- name: GetShows :many
SELECT *
FROM shows
LIMIT 10 OFFSET $1;