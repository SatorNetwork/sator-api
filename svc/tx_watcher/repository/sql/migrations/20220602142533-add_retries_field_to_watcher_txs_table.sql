-- +migrate Up
ALTER TABLE watcher_transactions
    ADD COLUMN retries INT NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE watcher_transactions
DROP COLUMN retries;
