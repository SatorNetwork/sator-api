-- +migrate Up
CREATE TABLE IF NOT EXISTS episode_access (
    episode_id uuid NOT NULL,
    user_id uuid NOT NULL,
    activated_at TIMESTAMP DEFAULT NULL
);
-- +migrate Down
DROP TABLE IF EXISTS episode_access;
