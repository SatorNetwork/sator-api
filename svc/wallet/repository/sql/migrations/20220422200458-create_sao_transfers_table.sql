-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION token_transfers_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS token_transfers (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    sender_address VARCHAR NOT NULL,
    recipient_address VARCHAR NOT NULL,
    tx_hash VARCHAR DEFAULT NULL,
    amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    status INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TRIGGER update_token_transfers_modtime BEFORE
UPDATE ON token_transfers FOR EACH ROW EXECUTE PROCEDURE token_transfers_update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_token_transfers_modtime ON token_transfers;
DROP TABLE IF EXISTS token_transfers;
DROP FUNCTION IF EXISTS token_transfers_update_updated_at_column();