-- +migrate Up
DROP TABLE IF EXISTS passed_challenges_data;
CREATE TABLE IF NOT EXISTS passed_challenges_data (
    user_id uuid NOT NULL,
    challenge_id uuid NOT NULL,
    reward_amount DOUBLE PRECISION NOT NULL DEFAULT 0
);
-- +migrate Down
DROP TABLE IF EXISTS passed_challenges_data;
CREATE TABLE IF NOT EXISTS passed_challenges_data (
        user_id uuid NOT NULL,
        challenge_id uuid NOT NULL,
        reward_amount DOUBLE PRECISION NOT NULL DEFAULT 0,
        PRIMARY KEY (user_id, challenge_id)
    );
