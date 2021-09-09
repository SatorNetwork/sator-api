-- +migrate Up
ALTER TABLE ratings DROP CONSTRAINT ratings_pkey;
ALTER TABLE ratings
ADD COLUMN id uuid NOT NULL DEFAULT uuid_generate_v4(),
ADD COLUMN title VARCHAR DEFAULT NULL,
ADD COLUMN review VARCHAR DEFAULT NULL,
ADD COLUMN username VARCHAR DEFAULT NULL;
ALTER TABLE ratings ADD PRIMARY KEY (id);
CREATE UNIQUE INDEX ratings_episode_user ON ratings USING BTREE (episode_id,user_id);
CREATE INDEX ratings_created_at ON ratings USING BTREE (created_at);

-- +migrate Down
ALTER TABLE ratings DROP CONSTRAINT ratings_pkey;
DROP INDEX IF EXISTS ratings_episode_user, ratings_created_at;
ALTER TABLE ratings 
DROP COLUMN IF EXISTS id,
DROP COLUMN IF EXISTS title,
DROP COLUMN IF EXISTS review,
DROP COLUMN IF EXISTS username;
ALTER TABLE ratings ADD PRIMARY KEY (episode_id,user_id);