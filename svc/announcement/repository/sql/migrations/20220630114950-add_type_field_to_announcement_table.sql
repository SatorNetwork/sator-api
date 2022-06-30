-- +migrate Up
ALTER TABLE announcements
    ADD COLUMN type VARCHAR NOT NULL DEFAULT '',
    ADD COLUMN type_specific_params VARCHAR NOT NULL DEFAULT '{}';

-- +migrate Down
ALTER TABLE announcements DROP COLUMN type, type_specific_params;