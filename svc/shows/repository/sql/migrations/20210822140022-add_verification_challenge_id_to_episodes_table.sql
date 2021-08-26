-- +migrate Up
ALTER TABLE episodes
ADD COLUMN verification_challenge_id uuid DEFAULT NULL;
-- +migrate Down
ALTER TABLE episodes DROP COLUMN verification_challenge_id;