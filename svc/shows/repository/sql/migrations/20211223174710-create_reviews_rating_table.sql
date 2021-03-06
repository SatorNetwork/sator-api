-- +migrate Up
CREATE TABLE IF NOT EXISTS reviews_rating (
    review_id uuid NOT NULL,
    user_id uuid NOT NULL,
    rating_type INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY(review_id) REFERENCES ratings(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(review_id, user_id)
);

-- +migrate Down
DROP TABLE IF EXISTS reviews_rating;

