-- name: GetEpisodes :many
SELECT *
FROM episodes
ORDER BY episode_number DESC
    LIMIT $1 OFFSET $2;
-- name: GetEpisodeByID :one
SELECT *
FROM episodes
WHERE id = $1;