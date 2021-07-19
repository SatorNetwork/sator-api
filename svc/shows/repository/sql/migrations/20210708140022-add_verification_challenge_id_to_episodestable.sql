-- +migrate Up
ALTER TABLE episodes
ADD COLUMN challenge_id uuid DEFAULT NULL;
-- +migrate Down
ALTER TABLE episodes DROP COLUMN challenge_id;