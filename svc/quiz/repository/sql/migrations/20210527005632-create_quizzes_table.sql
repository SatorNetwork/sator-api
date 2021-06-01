-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION quizzes_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS quizzes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    challenge_id uuid NOT NULL,
    prize_pool DOUBLE PRECISION NOT NULL,
    players_to_start INT NOT NULL,
    time_per_question INT NOT NULL DEFAULT 10,
    status INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX quizzes_challenge_id ON quizzes USING BTREE (challenge_id, status);
CREATE INDEX quizzes_created_at ON quizzes USING BTREE (created_at);
CREATE TRIGGER update_quizzes_modtime BEFORE
UPDATE ON quizzes FOR EACH ROW EXECUTE PROCEDURE quizzes_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_quizzes_modtime ON quizzes;
DROP TABLE IF EXISTS quizzes;
DROP FUNCTION IF EXISTS quizzes_update_updated_at_column();