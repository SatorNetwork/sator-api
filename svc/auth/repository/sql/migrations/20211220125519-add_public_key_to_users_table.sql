-- +migrate Up
ALTER TABLE users
    ADD COLUMN public_key VARCHAR DEFAULT NULL;
-- +migrate Down
ALTER TABLE users DROP COLUMN public_key;