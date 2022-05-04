-- name: GetPuzzleGameByID :one
SELECT *
FROM puzzle_games
WHERE id = $1;

-- name: GetPuzzleGameByEpisodeID :one
SELECT *
FROM puzzle_games
WHERE episode_id = $1;

-- name: CreatePuzzleGame :one
INSERT INTO puzzle_games (
    episode_id,
    prize_pool,
    parts_x,
    parts_y
)
VALUES (
    @episode_id,
    @prize_pool,
    @parts_x,
    @parts_x
) RETURNING *;

-- name: UpdatePuzzleGame :one
UPDATE puzzle_games
SET
    prize_pool = @prize_pool,
    parts_x = @parts_x,
    parts_y = @parts_x
WHERE id = @id
   RETURNING *;

-- name: LinkImageToPuzzleGame :exec
INSERT INTO puzzle_games_to_images (
    file_id,
    puzzle_game_id
)
VALUES (
    @file_id,
    @puzzle_game_id
);

-- name: GetPuzzleGameImageIDs :many
SELECT file_id FROM puzzle_games_to_images
WHERE puzzle_game_id = $1;

-- name: UnlinkImageFromPuzzleGame :exec
DELETE FROM puzzle_games_to_images
WHERE file_id = @file_id AND puzzle_game_id = @puzzle_game_id;

-- name: GetUserAvailableSteps :one
SELECT coalesce(steps, 0) FROM puzzle_games_attempts
WHERE user_id = $1 AND puzzle_game_id = $2 AND status = 0
ORDER BY created_at DESC
LIMIT 1;

-- name: GetPuzzleGameCurrentAttempt :one
SELECT *
FROM puzzle_games_attempts
WHERE user_id = $1 AND puzzle_game_id = $2 AND status = $3
ORDER BY created_at DESC
LIMIT 1;

-- name: UnlockPuzzleGame :one
INSERT INTO puzzle_games_attempts (
    puzzle_game_id,
    user_id,
    steps
) 
VALUES (
    @puzzle_game_id,
    @user_id,
    @steps
) RETURNING *;

-- name: StartPuzzleGame :one
UPDATE puzzle_games_attempts
SET status = 1, image = @image, tiles = @tiles
WHERE puzzle_game_id = @puzzle_game_id 
AND user_id = @user_id
AND status = 0
AND image IS NULL
RETURNING *;

-- name: FinishPuzzleGame :exec
UPDATE puzzle_games_attempts
SET status = 2,
    steps_taken = @steps_taken,
    rewards_amount = @rewards_amount,
    bonus_amount = @bonus_amount
WHERE puzzle_game_id = @puzzle_game_id
AND user_id = @user_id
AND status = 1
AND image IS NOT NULL
RETURNING *;

-- name: GetPuzzleGameUnlockOption :one
SELECT * FROM puzzle_game_unlock_options
WHERE id = $1;

-- name: GetPuzzleGameUnlockOptions :many
SELECT * FROM puzzle_game_unlock_options
WHERE disabled = FALSE
ORDER BY locked ASC, amount ASC;

-- name: UpdatePuzzleGameAttempt :one
UPDATE puzzle_games_attempts
SET
    status = @status,
    steps = @steps,
    steps_taken = @steps_taken,
    tiles = @tiles
WHERE id = @id
    RETURNING *;
