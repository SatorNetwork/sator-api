-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION solana_errors_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS solana_errors (
    provider_name VARCHAR NOT NULL,
    error_message VARCHAR NOT NULL,
    counter INT NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(provider_name, error_message)
    );
CREATE TRIGGER update_solana_errors_modtime BEFORE
    UPDATE ON solana_errors FOR EACH ROW EXECUTE PROCEDURE solana_errors_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_solana_errors_modtime ON solana_errors;
DROP TABLE IF EXISTS solana_errors;
DROP FUNCTION IF EXISTS solana_errors_update_updated_at_column();
