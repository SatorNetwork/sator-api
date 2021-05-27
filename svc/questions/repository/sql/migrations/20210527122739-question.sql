-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION questions_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS questions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    challenge_id uuid NOT NULL,
    question VARCHAR NOT NULL,
    question_order INT NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE INDEX ordering_questions_list ON questions USING BTREE (updated_at, created_at);
CREATE TRIGGER update_questions_modtime BEFORE
    UPDATE ON questions FOR EACH ROW EXECUTE PROCEDURE questions_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_questions_modtime ON questions;
DROP TABLE IF EXISTS questions;
DROP FUNCTION IF EXISTS questions_update_updated_at_column();