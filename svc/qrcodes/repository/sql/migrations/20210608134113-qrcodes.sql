-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION qrcodes_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS qrcodes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    show_id uuid NOT NULL,
    episode_id uuid NOT NULL,
    starts_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX qrcodes_created_at ON profiles USING BTREE (created_at);
CREATE TRIGGER update_qrcodes_modtime BEFORE
UPDATE ON qrcodes FOR EACH ROW EXECUTE PROCEDURE qrcodes_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_qrcodes_modtime ON qrcodes;
DROP TABLE IF EXISTS qrcodes;
DROP FUNCTION IF EXISTS qrcodes_update_updated_at_column();