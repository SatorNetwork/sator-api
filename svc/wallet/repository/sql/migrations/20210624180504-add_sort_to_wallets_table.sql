-- +migrate Up
ALTER TABLE wallets
ADD COLUMN sort INT NOT NULL DEFAULT 0;
-- +migrate Down
ALTER TABLE wallets DROP COLUMN sort;