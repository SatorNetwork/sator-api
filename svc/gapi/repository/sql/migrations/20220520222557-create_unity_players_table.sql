
-- +migrate Up
-- +migrate StatementBegin
CREATE
OR REPLACE FUNCTION unity_game_players_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS unity_game_players (
    user_id uuid PRIMARY KEY,
    energy_points INT NOT NULL DEFAULT 1,
    energy_refilled_at TIMESTAMP NOT NULL DEFAULT now(),
    selected_nft_id VARCHAR DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TRIGGER update_unity_game_players_modtime BEFORE
UPDATE ON unity_game_players FOR EACH ROW EXECUTE PROCEDURE unity_game_players_update_updated_at_column();

-- +migrate Down
DROP TRIGGER IF EXISTS update_unity_game_players_modtime ON unity_game_players;
DROP TABLE IF EXISTS unity_game_players;
DROP FUNCTION IF EXISTS unity_game_players_update_updated_at_column();
