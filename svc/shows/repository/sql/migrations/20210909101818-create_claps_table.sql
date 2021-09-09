-- +migrate Up
CREATE TABLE IF NOT EXISTS show_claps (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    show_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY(show_id) REFERENCES shows(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE IF EXISTS show_claps;
