-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION quiz_players_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS quiz_players (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    quiz_id uuid NOT NULL,
    user_id uuid NOT NULL,
    status INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX quiz_players_challenge_id ON quiz_players USING BTREE (challenge_id, status);
CREATE INDEX quiz_players_created_at ON quiz_players USING BTREE (created_at);
CREATE TRIGGER update_quiz_players_modtime BEFORE
UPDATE ON quiz_players FOR EACH ROW EXECUTE PROCEDURE quiz_players_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_quiz_players_modtime ON quiz_players;
DROP TABLE IF EXISTS quiz_players;
DROP FUNCTION IF EXISTS quiz_players_update_updated_at_column();