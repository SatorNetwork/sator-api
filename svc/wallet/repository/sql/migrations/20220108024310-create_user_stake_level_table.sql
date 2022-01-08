-- +migrate Up
CREATE TABLE IF NOT EXISTS stake_levels (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    min_stake_amount INT DEFAULT 0,
    title VARCHAR NOT NULL,
    subtitle VARCHAR NOT NULL,
    multiplier INT DEFAULT 0,
    PRIMARY KEY(id)
);

-- +migrate Down
DROP TABLE IF EXISTS stake_levels;