-- +migrate Up
CREATE TYPE episodes_status_type AS ENUM (
    'draft',
    'published',
    'archived'
);

ALTER TABLE episodes
ADD COLUMN status episodes_status_type DEFAULT 'draft';

UPDATE episodes
SET status = 'published'
WHERE archived = false;

UPDATE episodes
SET status = 'archived'
WHERE archived = true;

ALTER TABLE episodes
DROP COLUMN archived;

-- +migrate Down
ALTER TABLE episodes ADD COLUMN archived BOOLEAN DEFAULT FALSE NOT NULL;

UPDATE episodes
SET archived = true
WHERE status = 'archived';

UPDATE episodes
SET archived = false
WHERE status = 'published';

ALTER TABLE episodes
DROP COLUMN status;

DROP TYPE IF EXISTS episodes_status_type;