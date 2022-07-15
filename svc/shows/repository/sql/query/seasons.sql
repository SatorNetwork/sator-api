-- name: GetSeasonsByShowID :many
SELECT *
FROM seasons
WHERE show_id = $1 AND seasons.deleted_at IS NULL
ORDER BY season_number DESC
LIMIT $2 OFFSET $3;

-- name: GetSeasonByID :one
SELECT *
FROM seasons
WHERE id = $1 AND seasons.deleted_at IS NULL;

-- name: AddSeason :one
INSERT INTO seasons (
    show_id,
    season_number
) VALUES (
    @show_id,
    @season_number
) RETURNING *;

-- name: DeleteSeasonByID :exec
UPDATE seasons
SET deleted_at = NOW()
WHERE id = @id AND seasons.deleted_at IS NULL;

-- name: DeleteSeasonByShowID :exec
UPDATE seasons
SET deleted_at = NOW()
WHERE show_id = @show_id AND seasons.deleted_at IS NULL;