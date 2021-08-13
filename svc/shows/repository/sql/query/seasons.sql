-- name: GetSeasonsByShowID :many
SELECT *
FROM seasons
WHERE show_id = $1
ORDER BY season_number DESC
LIMIT $2 OFFSET $3;
-- name: GetSeasonByID :one
SELECT *
FROM seasons
WHERE id = $1;
-- name: AddSeason :exec
INSERT INTO seasons (
    show_id,
    season_number
) VALUES (
    @show_id,
    @season_number
);
-- name: DeleteSeasonByID :exec
DELETE FROM seasons
WHERE id = @id;