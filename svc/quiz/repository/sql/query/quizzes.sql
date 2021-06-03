-- name: AddNewQuiz :one
INSERT INTO quizzes (
        challenge_id,
        prize_pool,
        players_to_start,
        time_per_question
    )
VALUES (
        @challenge_id,
        @prize_pool,
        @players_to_start,
        @time_per_question
    ) RETURNING *;
-- name: UpdateQuizStatus :exec
UPDATE quizzes
SET status = @status
WHERE id = @id;
-- name: GetQuizByID :one
SELECT *
FROM quizzes
WHERE id = $1;
-- name: GetQuizByChallengeID :one
SELECT *
FROM quizzes
WHERE challenge_id = $1
    AND status = 0
ORDER BY created_at DESC;