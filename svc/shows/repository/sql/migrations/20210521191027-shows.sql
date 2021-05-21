-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd
CREATE TABLE IF NOT EXISTS shows (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    cover VARCHAR NOT NULL,
    has_new_episode BIT NOT NULL,
);
-- +migrate Down
DROP TABLE IF EXISTS shows;