-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION question_options_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS question_options (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id uuid NOT NULL,
    question_option VARCHAR NOT NULL,
    is_correct BOOLEAN,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE INDEX question_options_created_at ON question_options USING BTREE (updated_at, created_at);
CREATE TRIGGER update_question_options_modtime BEFORE
    UPDATE ON question_options FOR EACH ROW EXECUTE PROCEDURE question_options_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_question_options_modtime ON question_options;
DROP TABLE IF EXISTS question_options;
DROP FUNCTION IF EXISTS question_options_update_updated_at_column();