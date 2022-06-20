-- +migrate Up
CREATE TABLE IF NOT EXISTS flags (
    key VARCHAR PRIMARY KEY,
    value VARCHAR NOT NULL
);
-- +migrate Down
DROP TABLE IF EXISTS flags;