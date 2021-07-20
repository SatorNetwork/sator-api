-- +migrate Up
ALTER TABLE shows DROP COLUMN category;
-- +migrate Down
ALTER TABLE shows
    ADD COLUMN category VARCHAR DEFAULT NULL;