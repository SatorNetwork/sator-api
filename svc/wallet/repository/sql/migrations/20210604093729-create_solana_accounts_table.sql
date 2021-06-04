-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION solana_accounts_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS solana_accounts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_type VARCHAR NOT NULL,
    public_key VARCHAR NOT NULL,
    private_key BYTEA NOT NULL,
    status INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX solana_accounts_address ON solana_accounts USING BTREE (public_key);
CREATE TRIGGER update_solana_accounts_modtime BEFORE
UPDATE ON solana_accounts FOR EACH ROW EXECUTE PROCEDURE solana_accounts_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_solana_accounts_modtime ON solana_accounts;
DROP TABLE IF EXISTS solana_accounts;
DROP FUNCTION IF EXISTS solana_accounts_update_updated_at_column();