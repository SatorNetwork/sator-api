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
    result = $2,
    electricity_costs = $3,
    finished_at = now()
WHERE id = $4;

-- name: GetCurrentGame :one
SELECT * FROM unity_game_results 
WHERE user_id = $1 AND finished_at IS NULL 
ORDER BY created_at DESC LIMIT 1;