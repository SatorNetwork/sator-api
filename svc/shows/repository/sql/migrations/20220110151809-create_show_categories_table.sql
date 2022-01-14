-- +migrate Up
CREATE TABLE IF NOT EXISTS show_categories (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    disabled BOOLEAN DEFAULT FALSE,
    sort INT NOT NULL DEFAULT 0
    );
CREATE TABLE IF NOT EXISTS shows_to_categories (
    category_id uuid NOT NULL,
    show_id uuid NOT NULL,
    PRIMARY KEY(category_id, show_id),
    FOREIGN KEY(category_id) REFERENCES show_categories(id) ON DELETE CASCADE,
    FOREIGN KEY(show_id) REFERENCES shows(id) ON DELETE CASCADE
    );
-- +migrate Down
DROP TABLE IF EXISTS show_categories;
DROP TABLE IF EXISTS shows_to_categories;