-- +migrate Up
ALTER TABLE unity_game_players
    ADD COLUMN electricity_spent INT NOT NULL DEFAULT 0,
    ADD COLUMN electricity_costs DOUBLE PRECISION NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE unity_game_players 
    DROP COLUMN electricity_spent, 
    DROP COLUMN electricity_costs;