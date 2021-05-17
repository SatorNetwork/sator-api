-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION profiles_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS profiles (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    first_name VARCHAR DEFAULT NULL,
    last_name VARCHAR DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX profiles_user_id ON profiles USING BTREE (user_id);
CREATE INDEX profiles_created_at ON profiles USING BTREE (created_at);
CREATE TRIGGER update_profiles_modtime BEFORE
UPDATE ON profiles FOR EACH ROW EXECUTE PROCEDURE profiles_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_profiles_modtime ON profiles;
DROP TABLE IF EXISTS profiles;
DROP FUNCTION IF EXISTS profiles_update_updated_at_column();