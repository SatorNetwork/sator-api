-- +migrate Up
DROP TABLE IF EXISTS attempts;
CREATE TABLE IF NOT EXISTS attempts ( 
    user_id uuid NOT NULL,
    episode_id uuid NOT NULL,
    question_id uuid NOT NULL,
    answer_id uuid DEFAULT NULL,
    valid BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX attempts_user_episode_question ON attempts USING BTREE (user_id,episode_id,question_id);
-- +migrate Down
DROP TABLE IF EXISTS attempts;
CREATE TABLE IF NOT EXISTS attempts (
        user_id uuid NOT NULL,
        episode_id uuid NOT NULL,
        question_id uuid NOT NULL,
        answer_id uuid DEFAULT NULL,
        valid BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT NOW()
    );