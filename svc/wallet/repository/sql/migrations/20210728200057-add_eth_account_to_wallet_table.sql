-- +migrate Up
ALTER TABLE wallets
    ADD COLUMN ethereum_account_id uuid;
-- +migrate Down
ALTER TABLE wallets DROP COLUMN ethereum_account_id;