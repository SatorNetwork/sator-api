-- +migrate Up
CREATE UNIQUE INDEX rewards_user_rel_uniq ON rewards (user_id, relation_id) WHERE (transaction_type = 1);

-- +migrate Down
DROP INDEX IF EXISTS rewards_user_rel_uniq;