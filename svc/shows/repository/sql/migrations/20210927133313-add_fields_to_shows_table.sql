-- +migrate Up
ALTER TABLE shows
    ADD COLUMN realms_title VARCHAR DEFAULT NULL,
    ADD COLUMN realms_subtitle VARCHAR DEFAULT NULL,
    ADD COLUMN watch VARCHAR DEFAULT NULL;
-- +migrate Down
ALTER TABLE shows
    DROP COLUMN realms_title,
    DROP COLUMN realms_subtitle,
    DROP COLUMN watch;