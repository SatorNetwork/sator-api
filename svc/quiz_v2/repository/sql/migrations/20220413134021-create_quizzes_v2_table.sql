-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION quizzes_v2_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS quizzes_v2 (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    challenge_id uuid NOT NULL,
    distributed_rewards DOUBLE PRECISION NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY(challenge_id) REFERENCES challenges(id) ON DELETE CASCADE
    );
CREATE TRIGGER update_quizzes_v2_modtime BEFORE
    UPDATE ON quizzes_v2 FOR EACH ROW EXECUTE PROCEDURE quizzes_v2_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_quizzes_v2_modtime ON quizzes_v2;
DROP TABLE IF EXISTS quizzes_v2;
DROP FUNCTION IF EXISTS quizzes_v2_update_updated_at_column();
