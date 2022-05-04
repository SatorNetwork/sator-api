-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION iap_products_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS iap_products (
    id VARCHAR PRIMARY KEY,
    price_in_sao DOUBLE PRECISION NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE TRIGGER update_iap_products_modtime BEFORE
    UPDATE ON iap_products FOR EACH ROW EXECUTE PROCEDURE iap_products_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_iap_products_modtime ON iap_products;
DROP TABLE IF EXISTS iap_products;
DROP FUNCTION IF EXISTS iap_products_update_updated_at_column();
