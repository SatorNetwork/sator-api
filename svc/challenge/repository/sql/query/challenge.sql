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
-- name: AddChallenge :one
INSERT INTO challenges (
    show_id,
    title,
    description,
    prize_pool,
    players_to_start,
    time_per_question,
    updated_at,
    episode_id,
    kind,
    user_max_attempts
)
VALUES (
           @show_id,
           @title,
           @description,
           @prize_pool,
           @players_to_start,
           @time_per_question,
           @updated_at,
           @episode_id,
           @kind,
           @user_max_attempts
       ) RETURNING *;
-- name: UpdateChallenge :exec
UPDATE challenges
SET show_id = @show_id,
    title = @title,
    description = @description,
    prize_pool = @prize_pool,
    players_to_start = @players_to_start,
    time_per_question = @time_per_question,
    updated_at = @updated_at,
    episode_id = @episode_id,
    kind = @kind,
    user_max_attempts = @user_max_attempts
WHERE id = @id;
-- name: DeleteChallengeByID :exec
DELETE FROM challenges
WHERE id = @id;
-- name: GetChallengeByEpisodeID :one
SELECT *
FROM challenges
WHERE episode_id = $1
ORDER BY created_at DESC
    LIMIT 1;
-- name: GetChallengeByTitle :one
SELECT *
FROM challenges
WHERE title = $1
ORDER BY created_at DESC
    LIMIT 1;