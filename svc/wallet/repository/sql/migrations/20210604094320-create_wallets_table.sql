-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION wallets_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS wallets (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    solana_account_id uuid NOT NULL,
    wallet_name VARCHAR NOT NULL,
    status INT DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX wallets_user_id ON wallets USING BTREE (user_id, solana_account_id);
CREATE TRIGGER update_wallets_modtime BEFORE
UPDATE ON wallets FOR EACH ROW EXECUTE PROCEDURE wallets_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_wallets_modtime ON wallets;
DROP TABLE IF EXISTS wallets;
DROP FUNCTION IF EXISTS wallets_update_updated_at_column();