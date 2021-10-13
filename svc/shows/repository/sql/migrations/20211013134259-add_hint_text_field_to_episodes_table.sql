-- +migrate Up
ALTER TABLE episodes
ADD COLUMN hint_text VARCHAR DEFAULT NULL;
-- +migrate Down
ALTER TABLE episodes DROP COLUMN hint_text;