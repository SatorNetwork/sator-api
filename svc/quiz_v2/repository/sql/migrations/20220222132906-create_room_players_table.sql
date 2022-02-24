-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION room_players_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS room_players (
    challenge_id uuid NOT NULL,
    user_id uuid NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(challenge_id, user_id)
);
CREATE TRIGGER update_room_players_modtime BEFORE
    UPDATE ON room_players FOR EACH ROW EXECUTE PROCEDURE room_players_update_updated_at_column();
-- +migrate Down
DROP TRIGGER IF EXISTS update_room_players_modtime ON room_players;
DROP TABLE IF EXISTS room_players;
DROP FUNCTION IF EXISTS room_players_update_updated_at_column();
