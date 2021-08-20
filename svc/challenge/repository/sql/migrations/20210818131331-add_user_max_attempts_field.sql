-- +migrate Up
ALTER TABLE challenges
    ADD COLUMN user_max_attempts INT NOT NULL DEFAULT 2;
-- +migrate Down
ALTER TABLE challenges DROP COLUMN user_max_attempts;