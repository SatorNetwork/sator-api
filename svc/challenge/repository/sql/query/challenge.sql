-- name: GetChallenges :many
SELECT *
FROM challenges
WHERE show_id = $1
ORDER BY updated_at DESC,
    created_at DESC
LIMIT $2 OFFSET $3;
-- name: GetChallengeByID :one
SELECT *
FROM challenges
WHERE id = $1
ORDER BY created_at DESC
LIMIT 1;
-- name: GetChallengeByEpisodeID :one
SELECT *
FROM challenges
WHERE episode_id = $1
ORDER BY created_at DESC
LIMIT 1;