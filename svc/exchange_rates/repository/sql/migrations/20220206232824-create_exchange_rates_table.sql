-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION exchange_rates_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS exchange_rates (
    asset_type VARCHAR PRIMARY KEY,
    usd_price FLOAT NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE TRIGGER update_exchange_rates_modtime BEFORE
    UPDATE ON exchange_rates FOR EACH ROW EXECUTE PROCEDURE exchange_rates_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_exchange_rates_modtime ON exchange_rates;
DROP TABLE IF EXISTS exchange_rates;
DROP FUNCTION IF EXISTS exchange_rates_update_updated_at_column();
