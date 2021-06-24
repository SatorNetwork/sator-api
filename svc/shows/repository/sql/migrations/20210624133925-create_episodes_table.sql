-- +migrate Up
CREATE TABLE IF NOT EXISTS episodes (
    id uuid NOT NULL,
    show_id uuid NOT NULL,
    episode_number INT NOT NULL DEFAULT 0,
    cover VARCHAR NOT NULL,
    title VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    release_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(id, show_id)
);
-- +migrate Down
DROP TABLE IF EXISTS episodes;