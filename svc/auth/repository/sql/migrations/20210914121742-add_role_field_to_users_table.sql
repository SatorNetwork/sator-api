-- +migrate Up
ALTER TABLE users
    ADD COLUMN role VARCHAR DEFAULT NULL;
-- +migrate Down
ALTER TABLE users DROP COLUMN role;