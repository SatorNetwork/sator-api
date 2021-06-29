-- name: GetEpisodesByShowID :many
SELECT *
FROM episodes
WHERE show_id = $1
ORDER BY episode_number DESC
    LIMIT $2 OFFSET $3;
-- name: GetEpisodeByID :one
SELECT *
FROM episodes
WHERE id = $1;
