-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION firebase_registration_tokens_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS firebase_registration_tokens (
    device_id VARCHAR PRIMARY KEY,
    user_id uuid NOT NULL,
    registration_token VARCHAR NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users(id)
    );
CREATE TRIGGER update_firebase_registration_tokens_modtime BEFORE
    UPDATE ON firebase_registration_tokens FOR EACH ROW EXECUTE PROCEDURE firebase_registration_tokens_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_firebase_registration_tokens_modtime ON firebase_registration_tokens;
DROP TABLE IF EXISTS firebase_registration_tokens;
DROP FUNCTION IF EXISTS firebase_registration_tokens_update_updated_at_column();
