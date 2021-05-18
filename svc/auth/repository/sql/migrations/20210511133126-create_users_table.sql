-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION users_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    password bytea DEFAULT NULL,
    disabled bool NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email ON users USING BTREE (email);
CREATE UNIQUE INDEX users_username ON users USING BTREE (username);
CREATE INDEX users_created_at ON users USING BTREE (created_at);
CREATE TRIGGER update_users_modtime BEFORE
UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE users_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_users_modtime ON users;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS users_update_updated_at_column();