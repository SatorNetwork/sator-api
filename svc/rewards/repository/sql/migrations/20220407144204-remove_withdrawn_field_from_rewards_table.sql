-- +migrate Up
ALTER TABLE rewards
DROP COLUMN withdrawn;
-- +migrate Down
ALTER TABLE rewards
    ADD COLUMN withdrawn VARCHAR DEFAULT FALSE;