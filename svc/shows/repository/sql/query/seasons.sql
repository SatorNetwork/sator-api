-- name: GetSeasonsByShowID :many
SELECT *
FROM seasons
WHERE show_id = $1
ORDER BY season_number DESC
LIMIT $2 OFFSET $3;