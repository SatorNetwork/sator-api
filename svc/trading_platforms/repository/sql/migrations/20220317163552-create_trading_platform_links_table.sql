-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION trading_platform_links_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS trading_platform_links (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    link  VARCHAR NOT NULL,
    logo  VARCHAR NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE TRIGGER update_trading_platform_links_modtime BEFORE
    UPDATE ON trading_platform_links FOR EACH ROW EXECUTE PROCEDURE trading_platform_links_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_trading_platform_links_modtime ON trading_platform_links;
DROP TABLE IF EXISTS trading_platform_links;
DROP FUNCTION IF EXISTS trading_platform_links_update_updated_at_column();
