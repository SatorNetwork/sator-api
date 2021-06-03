-- +migrate Up
CREATE TABLE IF NOT EXISTS quiz_answers (
    quiz_id uuid NOT NULL,
    user_id uuid NOT NULL,
    question_id uuid NOT NULL,
    answer_id uuid NOT NULL,
    is_correct BOOLEAN NOT NULL,
    rate INT NOT NULL DEFAULT 0,
    pts INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(quiz_id, user_id, question_id)
);
-- +migrate Down
DROP TABLE IF EXISTS quiz_answers;