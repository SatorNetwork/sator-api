-- name: StoreAnswer :one
INSERT INTO quiz_answers (
        quiz_id,
        user_id,
        question_id,
        answer_id,
        is_correct,
        rate,
        pts
    )
VALUES (
        @quiz_id,
        @user_id,
        @question_id,
        @answer_id,
        @is_correct,
        @rate,
        CASE
            WHEN @is_correct THEN COALESCE(
                (
                    SELECT CASE
                            WHEN COUNT(*) > 0 THEN 0
                        END AS pts
                    FROM quiz_answers
                    WHERE question_id = @question_id
                        AND quiz_id = @quiz_id
                        AND is_correct = TRUE
                    GROUP BY question_id
                ),
                2
            )
            ELSE 0
        END
    ) ON CONFLICT (quiz_id, user_id, question_id) DO 
UPDATE SET 
    answer_id = EXCLUDED.answer_id, 
    is_correct = EXCLUDED.is_correct, 
    pts = CASE
            WHEN EXCLUDED.is_correct THEN COALESCE(
                (
                    SELECT CASE
                            WHEN COUNT(*) > 0 THEN 0
                        END AS pts
                    FROM quiz_answers
                    WHERE question_id = EXCLUDED.question_id
                        AND quiz_id = EXCLUDED.quiz_id
                        AND is_correct = TRUE
                    GROUP BY question_id
                ),
                2
            )
            ELSE 0
        END
RETURNING *;

-- name: CountCorrectAnswers :one
SELECT COUNT(answer_id) AS correct_answers,
    COUNT(pts) AS pts
FROM quiz_answers
WHERE quiz_id = @quiz_id
    AND user_id = @user_id
    AND is_correct = TRUE
GROUP BY quiz_id
LIMIT 1;
-- name: GetQuizWinnners :many
SELECT quiz_players.quiz_id AS quiz_id,
    quiz_players.user_id AS user_id,
    quiz_players.username AS username,
    COUNT(quiz_answers.answer_id)::INT AS correct_answers,
    SUM(quiz_answers.rate)::INT AS rate,
    SUM(quiz_answers.pts)::INT AS pts
FROM quiz_answers
    JOIN quiz_players ON quiz_players.quiz_id = quiz_answers.quiz_id
    AND quiz_players.user_id = quiz_answers.user_id
WHERE quiz_answers.quiz_id = @quiz_id
    AND quiz_answers.is_correct = TRUE
GROUP BY quiz_players.user_id,
    quiz_players.quiz_id
HAVING COUNT(quiz_answers.answer_id)::INT = @correct_answers::INT;
-- name: GetAnswer :one
SELECT *
FROM quiz_answers
WHERE quiz_id = @quiz_id
    AND user_id = @user_id
    AND question_id = @question_id
LIMIT 1;