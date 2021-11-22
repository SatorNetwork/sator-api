-- +migrate Up
CREATE TABLE IF NOT EXISTS whitelist (
    allowed_type VARCHAR NOT NULL,
    allowed_value VARCHAR NOT NULL,
    PRIMARY KEY(allowed_type, allowed_value)
);

-- +migrate Down
DROP TABLE IF EXISTS whitelist;