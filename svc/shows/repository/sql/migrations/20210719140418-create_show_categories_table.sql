-- +migrate Up
CREATE TABLE IF NOT EXISTS shows_categories (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_name VARCHAR NOT NULL,
    title VARCHAR NOT NULL,
    disabled BOOLEAN DEFAULT FALSE
);
-- +migrate Down
DROP TABLE IF EXISTS show_categories;