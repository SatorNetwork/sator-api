-- +migrate Up
ALTER TABLE solana_accounts DROP COLUMN account_type;
-- +migrate Down
ALTER TABLE solana_accounts
    ADD COLUMN account_type VARCHAR NOT NULL; 