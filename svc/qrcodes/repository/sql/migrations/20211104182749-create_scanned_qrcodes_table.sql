-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION scanned_qrcodes_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS scanned_qrcodes (
    user_id uuid NOT NULL,
    qrcode_id uuid NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, qrcode_id)
    );
CREATE INDEX scanned_qrcodes_created_at ON scanned_qrcodes USING BTREE (created_at);
CREATE TRIGGER update_scanned_qrcodes_modtime BEFORE
    UPDATE ON scanned_qrcodes FOR EACH ROW EXECUTE PROCEDURE scanned_qrcodes_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_scanned_qrcodes_modtime ON scanned_qrcodes;
DROP TABLE IF EXISTS scanned_qrcodes;
DROP FUNCTION IF EXISTS scanned_qrcodes_update_updated_at_column();