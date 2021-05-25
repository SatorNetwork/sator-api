-- name: GetChallenges :many
SELECT *
FROM challenges WHERE show_id = $1
ORDER BY updated_at,
    created_at DESC
LIMIT $2 OFFSET $3;
-- name: GetChallengeByID :one
SELECT *
FROM challenges
WHERE id = $1
ORDER BY created_at DESC
LIMIT 1;