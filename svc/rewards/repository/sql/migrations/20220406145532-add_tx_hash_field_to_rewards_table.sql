-- +migrate Up
ALTER TABLE rewards
ADD COLUMN tx_hash VARCHAR DEFAULT NULL;
-- +migrate Down
ALTER TABLE rewards
DROP COLUMN tx_hash;