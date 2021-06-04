-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION transactions_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS transactions (
    sender_wallet_id uuid NOT NULL,
    recipient_wallet_id uuid NOT NULL,
    transaction_hash VARCHAR NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX transactions_created_at ON transactions USING BTREE (updated_at, created_at);
CREATE TRIGGER update_transactions_modtime BEFORE
UPDATE ON transactions FOR EACH ROW EXECUTE PROCEDURE transactions_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_transactions_modtime ON transactions;
DROP TABLE IF EXISTS transactions;
DROP FUNCTION IF EXISTS transactions_update_updated_at_column();