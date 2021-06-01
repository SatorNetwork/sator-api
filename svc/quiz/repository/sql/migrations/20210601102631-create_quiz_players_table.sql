-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION quiz_players_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS quiz_players (
    quiz_id uuid NOT NULL,
    user_id uuid NOT NULL,
    username VARCHAR NOT NULL,
    status INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(quiz_id, user_id)
);
CREATE INDEX quiz_players_quiz_id ON quiz_players USING BTREE (quiz_id, status);
CREATE INDEX quiz_players_created_at ON quiz_players USING BTREE (created_at);
CREATE TRIGGER update_quiz_players_modtime BEFORE
UPDATE ON quiz_players FOR EACH ROW EXECUTE PROCEDURE quiz_players_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_quiz_players_modtime ON quiz_players;
DROP TABLE IF EXISTS quiz_players;
DROP FUNCTION IF EXISTS quiz_players_update_updated_at_column();