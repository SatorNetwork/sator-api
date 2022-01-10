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
    unstake_date
)
VALUES (
           @user_id,
           @wallet_id,
           @stake_amount,
           @stake_duration,
           @unstake_date
       ) RETURNING *;
-- name: UpdateStake :exec
UPDATE stake
SET stake_amount = @stake_amount,
    stake_duration = @stake_duration,
    unstake_date = @unstake_date
WHERE user_id = @user_id;
-- name: DeleteStakeByUserID :exec
DELETE FROM stake
WHERE user_id = $1;
-- name: GetTotalStake :one
SELECT SUM(stake_amount)::DOUBLE PRECISION
FROM stake;
