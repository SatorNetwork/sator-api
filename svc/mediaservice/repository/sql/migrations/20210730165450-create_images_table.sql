-- +migrate Up
CREATE TABLE IF NOT EXISTS images (
    id uuid PRIMARY KEY,
    file_name VARCHAR NOT NULL,
    file_path VARCHAR NOT NULL,
    file_url VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE INDEX images_created_at ON images USING BTREE (created_at);
-- +migrate Down
DROP TABLE IF EXISTS images;