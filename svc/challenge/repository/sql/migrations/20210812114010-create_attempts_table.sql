-- +migrate Up
CREATE TABLE IF NOT EXISTS attempts (  -- TODO: rename to verification_questions_attempts
    user_id uuid NOT NULL,
    episode_id uuid NOT NULL,
    question_id uuid NOT NULL,
    answer_id uuid DEFAULT NULL,
    valid BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +migrate Down
DROP TABLE IF EXISTS attempts;
