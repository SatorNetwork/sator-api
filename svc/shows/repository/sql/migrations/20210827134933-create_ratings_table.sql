-- +migrate Up
CREATE TABLE IF NOT EXISTS ratings (
    episode_id uuid NOT NULL,
    user_id uuid NOT NULL,
    rating INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(episode_id, user_id),
    FOREIGN KEY(episode_id) REFERENCES episodes(id) ON DELETE CASCADE
    );
-- +migrate Down
DROP TABLE IF EXISTS ratings;
