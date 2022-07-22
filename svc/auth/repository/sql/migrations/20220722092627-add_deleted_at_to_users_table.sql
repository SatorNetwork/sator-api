-- +migrate Up
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP DEFAULT NULL;
CREATE INDEX users_deleted_at ON users USING BTREE (deleted_at);

-- +migrate Down
ALTER TABLE users DROP COLUMN deleted_at;