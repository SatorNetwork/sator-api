-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION ethereum_accounts_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS ethereum_accounts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    public_key BYTEA NOT NULL,
    private_key BYTEA NOT NULL,
    address VARCHAR NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE UNIQUE INDEX ethereum_accounts_address ON ethereum_accounts USING BTREE (public_key);
CREATE TRIGGER update_ethereum_accounts_modtime BEFORE
    UPDATE ON ethereum_accounts FOR EACH ROW EXECUTE PROCEDURE ethereum_accounts_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_ethereum_accounts_modtime ON ethereum_accounts;
DROP TABLE IF EXISTS ethereum_accounts;
DROP FUNCTION IF EXISTS ethereum_accounts_update_updated_at_column();