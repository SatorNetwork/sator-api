-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION answer_options_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS answer_options (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    question_id uuid NOT NULL,
    answer_option VARCHAR NOT NULL,
    is_correct BOOLEAN,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX answer_options_created_at ON answer_options USING BTREE (updated_at, created_at);
CREATE TRIGGER update_answer_options_modtime BEFORE
UPDATE ON answer_options FOR EACH ROW EXECUTE PROCEDURE answer_options_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_answer_options_modtime ON answer_options;
DROP TABLE IF EXISTS answer_options;
DROP FUNCTION IF EXISTS answer_options_update_updated_at_column();