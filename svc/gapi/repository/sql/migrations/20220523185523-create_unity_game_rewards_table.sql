
-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS unity_game_rewards (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    relation_id uuid DEFAULT NULL,
    operation_type int NOT NULL,
    amount DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS unity_game_rewards_user_id_idx ON unity_game_rewards (user_id);

-- +migrate Down
DROP TABLE IF EXISTS unity_game_rewards;
