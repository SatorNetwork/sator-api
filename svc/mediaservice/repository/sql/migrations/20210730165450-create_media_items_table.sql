-- +migrate Up
CREATE TABLE IF NOT EXISTS items (
    id uuid PRIMARY KEY,
    filename VARCHAR NOT NULL,
    filepath VARCHAR NOT NULL,
    relation_type VARCHAR,
    relation_id uuid,
    created_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE INDEX items_created_at ON items USING BTREE (created_at);
-- +migrate Down
DROP TABLE IF EXISTS items;