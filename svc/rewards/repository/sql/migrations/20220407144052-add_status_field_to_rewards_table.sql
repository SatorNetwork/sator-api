-- Status field :
--   0 - Available
--   1 - Requested
--   2 - Failed
--   3 - Withdrawn
-- +migrate Up
ALTER TABLE rewards
    ADD COLUMN status VARCHAR;
-- +migrate Down
ALTER TABLE rewards
DROP COLUMN status;
