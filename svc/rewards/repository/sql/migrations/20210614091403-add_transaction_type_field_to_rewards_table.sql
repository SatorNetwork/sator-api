-- +migrate Up
ALTER TABLE rewards
ADD COLUMN transaction_type INTEGER NOT NULL DEFAULT 1;
CREATE INDEX rewards_transaction_type ON rewards USING BTREE (transaction_type);
-- +migrate Down
DROP INDEX IF EXISTS rewards_transaction_type;
ALTER TABLE rewards DROP COLUMN transaction_type;