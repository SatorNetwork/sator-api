-- +migrate Up
CREATE TABLE IF NOT EXISTS episodes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    show_id uuid NOT NULL,
    episode_number INT NOT NULL DEFAULT 0,
    cover VARCHAR DEFAULT NULL,
    title VARCHAR NOT NULL,
    description VARCHAR DEFAULT NULL,
    release_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +migrate Down
DROP TABLE IF EXISTS episodes;