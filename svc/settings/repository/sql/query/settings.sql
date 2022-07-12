-- name: GetSettingByKey :one
SELECT *
FROM settings
WHERE key = $1;

-- name: AddSetting :one
INSERT INTO settings (key, name, value_type, value, description)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateSetting :one
UPDATE settings SET value = @value
WHERE key = @key
RETURNING *;

-- name: GetSettings :many
SELECT * FROM settings ORDER BY key;

-- name: DeleteSetting :exec
DELETE FROM settings WHERE key = $1;

