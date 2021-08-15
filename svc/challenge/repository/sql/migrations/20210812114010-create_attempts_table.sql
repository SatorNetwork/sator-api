-- +migrate Up
CREATE TABLE IF NOT EXISTS attempts (
    user_id uuid NOT NULL,
    episode_id uuid NOT NULL,
    question_id uuid NOT NULL,
    answer_id uuid DEFAULT NULL,
    valid BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, episode_id, question_id)
);
-- +migrate Down
DROP TABLE IF EXISTS attempts;
