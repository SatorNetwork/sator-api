-- +migrate Up
CREATE TYPE shows_status_type AS ENUM (
    'draft',
    'published',
    'archived'
);

ALTER TABLE shows
ADD COLUMN status shows_status_type DEFAULT 'draft';

UPDATE shows
SET status = 'published'
WHERE archived = false;

UPDATE shows
SET status = 'archived'
WHERE archived = true;

ALTER TABLE shows
DROP COLUMN archived;

-- +migrate Down
ALTER TABLE shows ADD COLUMN archived BOOLEAN DEFAULT FALSE NOT NULL;

UPDATE shows
SET archived = true
WHERE status = 'archived';

UPDATE shows
SET archived = false
WHERE status = 'published';

ALTER TABLE shows
DROP COLUMN status;

DROP TYPE IF EXISTS shows_status_type;