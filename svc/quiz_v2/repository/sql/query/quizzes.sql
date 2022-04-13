-- name: RegisterNewQuiz :one
INSERT INTO quizzes_v2 (
    challenge_id,
    distributed_rewards
)
VALUES (
    @challenge_id,
    @distributed_rewards
) RETURNING *;

-- name: GetDistributedRewardsByChallengeID :one
SELECT SUM(distributed_rewards)::DOUBLE PRECISION
FROM quizzes_v2
WHERE challenge_id = $1;

-- name: CleanUp :exec
DELETE FROM quizzes_v2;
