-- +migrate Up
CREATE TABLE IF NOT EXISTS episode_access (
    episode_id uuid NOT NULL,
    user_id uuid NOT NULL,
    activated_at TIMESTAMP DEFAULT NULL
);
CREATE INDEX episode_access_episode_user ON episode_access USING BTREE (episode_id,user_id);
CREATE INDEX episode_access_activated_at ON episode_access USING BTREE (activated_at);
-- +migrate Down
DROP TABLE IF EXISTS episode_access;
