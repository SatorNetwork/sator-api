-- +migrate Up
ALTER TABLE shows ADD COLUMN archived BOOLEAN DEFAULT FALSE NOT NULL;
ALTER TABLE episodes ADD COLUMN archived BOOLEAN DEFAULT FALSE NOT NULL;

-- +migrate Down
ALTER TABLE shows DROP COLUMN archived;
ALTER TABLE episodes DROP COLUMN archived;