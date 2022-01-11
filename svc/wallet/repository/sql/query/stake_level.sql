-- name: AddStakeLevel :one
INSERT INTO stake_levels (
        min_stake_amount,
        title,
        subtitle,
        multiplier
    )
VALUES (
        @min_stake_amount,
        @title,
        @subtitle,
        @multiplier
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
SET min_stake_amount = $2, title = $3, subtitle = $4, multiplier = $5
WHERE id = $1;
