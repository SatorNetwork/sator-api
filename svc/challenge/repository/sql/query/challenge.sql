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
-- name: AddChallenge :exec
INSERT INTO challenges (
    show_id,
    title,
    description,
    prize_pool,
    players_to_start,
    time_per_question,
    updated_at
)
VALUES (
           @show_id,
           @title,
           @description,
           @prize_pool,
           @players_to_start,
           @time_per_question,
           @updated_at
       );
-- name: UpdateChallenge :exec
UPDATE challenges
SET show_id = @show_id,
    title = @title,
    description = @description,
    prize_pool = @prize_pool,
    players_to_start = @players_to_start,
    time_per_question = @time_per_question,
    updated_at = @updated_at
WHERE id = @id;
-- name: DeleteChallengeByID :exec
DELETE FROM challenges
WHERE id = @id;