-- name: AddNewPlayer :one
INSERT INTO unity_game_players (user_id, energy_points, energy_refilled_at, selected_nft_id) 
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetPlayer :one
SELECT * FROM unity_game_players WHERE user_id = $1;

-- name: RefillEnergyOfPlayer :exec
UPDATE unity_game_players SET energy_points = $1, energy_refilled_at = $2 WHERE user_id = $3;

-- name: SpendEnergyOfPlayer :exec
UPDATE unity_game_players SET energy_points = energy_points-1 WHERE user_id = $1;

-- name: ResetEnergyRefilledAtOfPlayer :exec
UPDATE unity_game_players SET energy_refilled_at = now() WHERE user_id = $1;

-- name: StoreSelectedNFT :exec
UPDATE unity_game_players SET selected_nft_id = $1 WHERE user_id = $2;

-- name: AddElectricityToPlayer :exec
UPDATE unity_game_players 
SET 
    electricity_costs = electricity_costs + $1,
    electricity_spent = electricity_spent + $2 
WHERE user_id = $3;

-- name: ResetElectricityForPlayer :exec
UPDATE unity_game_players SET electricity_costs = 0, electricity_spent = 0 WHERE user_id = $1;

-- name: SelectCharacterToPlayer :exec
UPDATE unity_game_players
SET selected_character_id = $1
WHERE user_id = $2;
