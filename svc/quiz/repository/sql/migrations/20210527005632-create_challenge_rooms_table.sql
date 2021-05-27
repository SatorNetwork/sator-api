-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS challenge_rooms (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    challenge_id uuid NOT NULL,
    prize_pool DOUBLE PRECISION NOT NULL,
    players_to_start INT NOT NULL,
    time_per_question INT NOT NULL DEFAULT 10,
    status INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX challenge_rooms_challenge_id ON challenge_rooms USING BTREE (challenge_id, status);
CREATE INDEX challenge_rooms_created_at ON challenge_rooms USING BTREE (created_at);
-- +migrate Down
DROP TABLE IF EXISTS challenge_rooms;