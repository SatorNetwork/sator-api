-- Status field :
--   0 - Available
--   1 - Requested
--   2 - In progress
--   3 - Failed
--   4 - Done
-- +migrate Up
ALTER TABLE rewards
    ADD COLUMN status INTEGER NOT NULL DEFAULT 0;
-- +migrate Down
ALTER TABLE rewards
DROP COLUMN status;