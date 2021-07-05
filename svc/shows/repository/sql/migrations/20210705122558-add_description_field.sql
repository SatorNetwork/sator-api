-- +migrate Up
ALTER TABLE shows
    ADD COLUMN description VARCHAR DEFAULT NULL;
-- +migrate Down
ALTER TABLE shows DROP COLUMN description;