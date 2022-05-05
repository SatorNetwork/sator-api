-- Status field :
--   0 - Available
--   1 - Requested
--   2 - Failed
--   3 - Done
-- +migrate Up
ALTER TABLE rewards
    ADD COLUMN status VARCHAR NOT NULL DEFAULT 'TransactionStatusAvailable';
-- +migrate Down
ALTER TABLE rewards
DROP COLUMN status;
