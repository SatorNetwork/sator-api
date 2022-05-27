-- name: AddNewPlayer :one
INSERT INTO unity_game_players (user_id, energy_points, energy_refilled_at, selected_nft_id) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetPlayer :one
SELECT * FROM unity_game_players WHERE user_id = $1;

-- name: RefillEnergyOfPlayer :exec
UPDATE unity_game_players SET energy_points = $1, energy_refilled_at = $2 WHERE user_id = $3;

-- name: SpendEnergyOfPlayer :exec
UPDATE unity_game_players SET energy_points = energy_points-1 WHERE user_id = $1;

-- name: StoreSelectedNFT :exec
UPDATE unity_game_players SET selected_nft_id = $1 WHERE user_id = $2;

