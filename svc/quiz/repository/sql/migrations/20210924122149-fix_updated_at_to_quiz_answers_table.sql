-- +migrate Up
DROP TRIGGER IF EXISTS update_quiz_answers_modtime ON quiz_answers;
ALTER TABLE quiz_answers DROP COLUMN updated_at;
ALTER TABLE quiz_answers
    ADD COLUMN updated_at TIMESTAMP DEFAULT NULL;
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION quiz_answers_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
CREATE TRIGGER update_quiz_answers_modtime BEFORE
    UPDATE ON quiz_answers FOR EACH ROW EXECUTE PROCEDURE quiz_answers_update_updated_at_column();
-- +migrate StatementEnd

-- +migrate Down
DROP TRIGGER IF EXISTS update_quiz_answers_modtime ON quiz_answers;
ALTER TABLE quiz_answers DROP COLUMN updated_at;