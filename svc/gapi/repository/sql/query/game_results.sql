-- name: StartGame :exec
INSERT INTO unity_game_results (
    user_id,
    nft_id,
    complexity,
    is_training
) VALUES ($1, $2, $3, $4);

-- name: FinishGame :exec
UPDATE unity_game_results SET
    blocks_done = $1,
    rewards = $2,
    finished_at = now()
WHERE id = $3;

-- name: GetCurrentGame :one
SELECT * FROM unity_game_results 
WHERE user_id = $1 AND finished_at IS NULL 
ORDER BY created_at DESC LIMIT 1;