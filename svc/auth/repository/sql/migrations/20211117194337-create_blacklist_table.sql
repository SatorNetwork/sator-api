-- +migrate Up
CREATE TABLE IF NOT EXISTS blacklist (
    restricted_type VARCHAR NOT NULL,
    restricted_value VARCHAR NOT NULL,
    PRIMARY KEY(restricted_type, restricted_value)
);

-- +migrate Down
DROP TABLE IF EXISTS blacklist;