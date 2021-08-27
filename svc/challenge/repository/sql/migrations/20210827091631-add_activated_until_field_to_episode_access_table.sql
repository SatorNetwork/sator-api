-- +migrate Up
ALTER TABLE episode_access
    ADD COLUMN activated_before TIMESTAMP DEFAULT NULL;
CREATE INDEX episode_access_activated_before ON episode_access USING BTREE (activated_before);
-- +migrate Down
DROP INDEX IF EXISTS episode_access_activated_before;
ALTER TABLE episode_access DROP COLUMN activated_before;