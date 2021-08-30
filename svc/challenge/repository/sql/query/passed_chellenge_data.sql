-- name: AddChallengeAttempt :one
INSERT INTO passed_challenges_data (user_id, challenge_id)
VALUES ($1, $2) RETURNING *;
-- name: StoreChallengeReceivedRewardAmount :exec
UPDATE passed_challenges_data
SET reward_amount = @reward_amount
WHERE user_id = @user_id AND challenge_id = @challenge_id;
-- name: CountPassedChallengeAttempts :one
SELECT COUNT (*)
FROM passed_challenges_data
WHERE user_id = $1 AND challenge_id = $2;
-- name: GetChallengeReceivedRewardAmount :one
SELECT reward_amount
FROM passed_challenges_data
WHERE user_id = $1 AND challenge_id = $2;
