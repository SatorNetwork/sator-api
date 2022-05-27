-- +migrate Up
ALTER TABLE unity_game_results
    ADD COLUMN result INT DEFAULT NULL,
    ADD COLUMN electricity_costs DOUBLE PRECISION NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE unity_game_results 
    DROP COLUMN result, 
    DROP COLUMN electricity_costs;