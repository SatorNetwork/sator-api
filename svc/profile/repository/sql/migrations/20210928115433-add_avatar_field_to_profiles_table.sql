-- +migrate Up
ALTER TABLE profiles
    ADD COLUMN avatar VARCHAR DEFAULT NULL;
-- +migrate Down
ALTER TABLE profiles DROP COLUMN avatar;