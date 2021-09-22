-- +migrate Up
ALTER TABLE users
    ADD COLUMN role VARCHAR NOT NULL DEFAULT "user";
-- +migrate Down
ALTER TABLE users DROP COLUMN role;