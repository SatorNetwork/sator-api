-- name: GetStakeByUserID :one
SELECT *
FROM stake
WHERE user_id = $1
    LIMIT 1;
-- name: AddStake :one
INSERT INTO stake (
    user_id,
    wallet_id,
    stake_amount,
    stake_duration,
    unstake_date,
    unstake_timestamp
)
VALUES (
    @user_id,
    @wallet_id,
    @stake_amount,
    @stake_duration,
    @unstake_date,
    @unstake_timestamp
) RETURNING *;
-- name: UpdateStake :exec
UPDATE stake
SET stake_amount = @stake_amount,
    stake_duration = @stake_duration,
    unstake_date = @unstake_date,
    unstake_timestamp = @unstake_timestamp
WHERE user_id = @user_id;
-- name: DeleteStakeByUserID :exec
DELETE FROM stake
WHERE user_id = $1;
-- name: GetTotalStake :one
SELECT coalesce(SUM(coalesce(stake_amount, 0)), 0)::DOUBLE PRECISION
FROM stake;
