-- +migrate Up
ALTER TABLE shows ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;
CREATE INDEX shows_deleted_at_idx ON shows USING BTREE (deleted_at);

ALTER TABLE episodes ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;
CREATE INDEX episodes_deleted_at_idx ON episodes USING BTREE (deleted_at);

-- +migrate Down
ALTER TABLE shows DROP COLUMN deleted_at;
ALTER TABLE episodes DROP COLUMN deleted_at;