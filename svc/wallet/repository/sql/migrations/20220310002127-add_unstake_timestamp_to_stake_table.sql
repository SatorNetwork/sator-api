-- +migrate Up
ALTER TABLE stake
    ADD COLUMN unstake_timestamp BIGINT NOT NULL DEFAULT 0;
-- +migrate Down
ALTER TABLE stake DROP COLUMN unstake_timestamp;
