-- +migrate Up
ALTER TABLE seasons ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;
CREATE INDEX seasons_deleted_at_idx ON seasons USING BTREE (deleted_at);

-- +migrate Down
ALTER TABLE seasons DROP COLUMN deleted_at;