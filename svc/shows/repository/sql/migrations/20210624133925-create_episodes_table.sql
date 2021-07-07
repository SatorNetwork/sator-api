-- +migrate Up
CREATE TABLE IF NOT EXISTS seasons (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    show_id uuid NOT NULL,
    season_number INT NOT NULL DEFAULT 0,
    FOREIGN KEY(show_id) REFERENCES shows(id) ON
    DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS episodes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    show_id uuid NOT NULL,
    season_id uuid DEFAULT NULL,
    episode_number INT NOT NULL DEFAULT 0,
    cover VARCHAR DEFAULT NULL,
    title VARCHAR NOT NULL,
    description VARCHAR DEFAULT NULL,
    release_date TIMESTAMP DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY(season_id) REFERENCES seasons(id) ON
    DELETE CASCADE
);
-- +migrate Down
DROP TABLE IF EXISTS episodes;
DROP TABLE IF EXISTS seasons;