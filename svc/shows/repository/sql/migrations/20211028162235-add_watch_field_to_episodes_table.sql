-- +migrate Up
ALTER TABLE episodes
    ADD COLUMN watch VARCHAR DEFAULT NULL;
-- +migrate Down
ALTER TABLE episodes DROP COLUMN watch;