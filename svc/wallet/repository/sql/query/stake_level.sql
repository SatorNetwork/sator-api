-- name: AddStakeLevel :one
INSERT INTO stake_levels (
        min_stake_amount,
        min_days_amount,
        title,
        subtitle,
        multiplier,
        disabled
    )
VALUES (
        @min_stake_amount,
        @min_days_amount,
        @title,
        @subtitle,
        @multiplier,
        @disabled
    ) ON CONFLICT (title) DO NOTHING RETURNING *;
-- name: GetStakeLevelByID :one
SELECT *
FROM stake_levels
WHERE id = @id
LIMIT 1;
-- name: GetAllStakeLevels :many
SELECT *
FROM stake_levels
ORDER BY min_stake_amount DESC;
-- name: UpdateStakeLevel :exec
UPDATE stake_levels
SET min_stake_amount = $2, min_days_amount= $3, title = $4, subtitle = $5, multiplier = $6, disabled = $7
WHERE id = $1;
