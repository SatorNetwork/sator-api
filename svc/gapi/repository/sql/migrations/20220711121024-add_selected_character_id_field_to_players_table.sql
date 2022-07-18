-- +migrate Up
ALTER TABLE unity_game_players
ADD COLUMN selected_character_id varchar;

-- +migrate Down
ALTER TABLE unity_game_players
DROP COLUMN selected_character_id;