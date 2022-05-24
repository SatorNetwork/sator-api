-- name: GetSettings :many
SELECT * FROM unity_game_settings ORDER BY key;

-- name: GetSettingByKey :one
SELECT * FROM unity_game_settings WHERE key = $1;

-- name: AddSetting :one
INSERT INTO unity_game_settings (key, name, value_type, value, description) 
VALUES ($1, $2, $3, $4, $5) 
RETURNING *;

-- name: UpdateSetting :one
UPDATE unity_game_settings 
SET value = $2 
WHERE key = $1 
RETURNING *;

-- name: DeleteSetting :exec
DELETE FROM unity_game_settings WHERE key = $1;