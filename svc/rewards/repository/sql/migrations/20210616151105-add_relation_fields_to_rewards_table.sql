-- +migrate Up
ALTER TABLE rewards
    RENAME COLUMN quiz_id TO relation_id;
ALTER TABLE rewards
    ADD COLUMN relation_type VARCHAR DEFAULT NULL;
CREATE INDEX rewards_relations ON rewards USING BTREE (relation_id, relation_type);
-- +migrate Down
ALTER TABLE rewards
    RENAME COLUMN relation_id TO quiz_id;
DROP INDEX IF EXISTS rewards_relations;
ALTER TABLE rewards DROP COLUMN relation_type;