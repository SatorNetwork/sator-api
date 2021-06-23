-- +migrate Up
ALTER TABLE shows
    ADD COLUMN category VARCHAR DEFAULT NULL;
-- +migrate Down
ALTER TABLE shows DROP COLUMN category;