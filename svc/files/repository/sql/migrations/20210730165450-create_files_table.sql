-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS files (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_name VARCHAR NOT NULL,
    file_path VARCHAR NOT NULL,
    file_url VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX files_created_at ON files USING BTREE (created_at);
-- +migrate Down
DROP TABLE IF EXISTS files;