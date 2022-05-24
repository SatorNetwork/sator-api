-- name: GetUserRewards :one 
WITH deposited AS (
    SELECT user_id,  SUM(amount)::DOUBLE PRECISION AS amount
    FROM unity_game_rewards
    WHERE operation_type = 1
    GROUP BY user_id
), withdrawn AS (
    SELECT user_id,  SUM(amount)::DOUBLE PRECISION AS amount
    FROM unity_game_rewards
    WHERE operation_type = 2
    GROUP BY user_id
)
SELECT (deposited.amount - withdrawn.amount)::DOUBLE PRECISION AS total_reward_amount
FROM unity_game_players
LEFT JOIN deposited ON unity_game_players.user_id = deposited.user_id
LEFT JOIN withdrawn ON unity_game_players.user_id = withdrawn.user_id
WHERE unity_game_players.user_id = @user_id;

-- name: GetUserRewardsDeposited :one
SELECT SUM(amount)::DOUBLE PRECISION AS total_reward_amount
FROM unity_game_rewards
WHERE user_id = @user_id
AND operation_type = 1;

-- name: GetUserRewardsWithdrawn :one
SELECT SUM(amount)::DOUBLE PRECISION AS total_reward_amount
FROM unity_game_rewards
WHERE user_id = @user_id
AND operation_type = 2;

-- name: RewardsDeposit :exec
INSERT INTO unity_game_rewards (user_id, relation_id, operation_type, amount)
VALUES (@user_id, @relation_id, 1, @amount);

-- name: RewardsWithdraw :exec
INSERT INTO unity_game_rewards (user_id, operation_type, amount)
VALUES (@user_id, 2, @amount);