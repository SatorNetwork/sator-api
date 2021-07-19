-- +migrate Up
CREATE TABLE IF NOT EXISTS shows_to_category (
    category_id uuid NOT NULL,
    show_id uuid NOT NULL,
    PRIMARY KEY(category_id, show_id),
    FOREIGN KEY(category_id) REFERENCES shows_categories(id) ON DELETE CASCADE,
    FOREIGN KEY(show_id) REFERENCES shows(id) ON DELETE CASCADE
);
-- +migrate Down
DROP TABLE IF EXISTS shows_to_category;