-- +migrate Up
ALTER TABLE challenges
ADD COLUMN episode_id uuid DEFAULT NULL;
-- +migrate Down
ALTER TABLE challenges DROP COLUMN episode_id;