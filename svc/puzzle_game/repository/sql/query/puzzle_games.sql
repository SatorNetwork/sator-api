-- name: GetPuzzleGameByID :one
SELECT *
FROM puzzle_games
WHERE id = $1;

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
    @parts_y
) RETURNING *;

-- name: UpdatePuzzleGame :one
UPDATE puzzle_games
SET
    episode_id = @episode_id,
    prize_pool = @prize_pool,
    parts_x = @parts_x,
    parts_y = @parts_y
WHERE id = @id
   RETURNING *;

-- name: AddImageToPuzzleGame :exec
INSERT INTO puzzle_games_to_images (
    file_id,
    puzzle_game_id
)
VALUES (
    @file_id,
    @puzzle_game_id
);

-- name: DeleteImageFromPuzzleGame :exec
DELETE FROM puzzle_games_to_images
WHERE file_id = @file_id AND puzzle_game_id = @puzzle_game_id;

-- -- name: DeleteTradingPlatformLink :exec
-- DELETE FROM trading_platform_links
-- WHERE id = @id;
--
-- -- name: GetTradingPlatformLinks :many
-- SELECT *
-- FROM trading_platform_links
-- ORDER BY title ASC
--    LIMIT $1 OFFSET $2;
