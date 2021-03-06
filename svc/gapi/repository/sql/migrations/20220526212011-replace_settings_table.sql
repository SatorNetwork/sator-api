
-- +migrate Up
DROP TABLE IF EXISTS unity_game_settings;
DROP TYPE IF EXISTS unity_game_settings_value_type;

CREATE TYPE unity_game_settings_value_type AS ENUM (
    'string',
    'int',
    'float',
    'bool',
    'json',
    'duration',
    'datetime'
);

CREATE TABLE IF NOT EXISTS unity_game_settings (
    key VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    value_type unity_game_settings_value_type NOT NULL,
    value text NOT NULL,
    description text DEFAULT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS unity_game_settings;
DROP TYPE IF EXISTS unity_game_settings_value_type;

CREATE TYPE unity_game_settings_value_type AS ENUM (
    'string',
    'int',
    'float',
    'bool',
    'json'
);

CREATE TABLE IF NOT EXISTS unity_game_settings (
    key VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    value_type unity_game_settings_value_type NOT NULL,
    value bytea NOT NULL,
    description text DEFAULT NULL
);