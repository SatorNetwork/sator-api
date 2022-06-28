-- +migrate Up
CREATE TYPE settings_value_type AS ENUM (
    'string',
    'int',
    'float',
    'bool',
    'json',
    'duration',
    'datetime'
);

CREATE TABLE IF NOT EXISTS settings (
   key VARCHAR PRIMARY KEY,
   name VARCHAR NOT NULL,
   value_type settings_value_type NOT NULL,
   value text NOT NULL,
   description text DEFAULT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS settings;
DROP TYPE IF EXISTS settings_value_type;
