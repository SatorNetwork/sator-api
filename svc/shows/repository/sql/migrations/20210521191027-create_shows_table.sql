-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION shows_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS shows (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    cover VARCHAR NOT NULL,
    has_new_episode BOOLEAN NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX ordering_shows_list ON shows USING BTREE (updated_at, created_at);
CREATE TRIGGER update_shows_modtime BEFORE
UPDATE ON shows FOR EACH ROW EXECUTE PROCEDURE shows_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_shows_modtime ON shows;
DROP TABLE IF EXISTS shows;
DROP FUNCTION IF EXISTS shows_update_updated_at_column();