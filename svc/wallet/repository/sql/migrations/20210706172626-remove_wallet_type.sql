-- +migrate Up
ALTER TABLE wallets DROP COLUMN wallet_type;
ALTER TABLE solana_accounts DROP COLUMN account_type;
-- +migrate Down
ALTER TABLE wallets
    ADD COLUMN wallet_type VARCHAR NOT NULL;
ALTER TABLE solana_accounts
    ADD COLUMN account_type VARCHAR NOT NULL;