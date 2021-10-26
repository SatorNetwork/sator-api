-- +migrate Up

-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE
OR REPLACE FUNCTION nft_update_updated_at_column() RETURNS TRIGGER AS $$
BEGIN NEW .updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd

CREATE TABLE IF NOT EXISTS nft_items (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id uuid DEFAULT NULL,
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    cover VARCHAR NOT NULL,
    supply BIGINT NOT NULL DEFAULT 1,
    buy_now_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    token_uri TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX ordering_nft_items_list ON nft_items USING BTREE (updated_at, created_at);
CREATE TRIGGER update_nft_items_modtime BEFORE
UPDATE ON nft_items FOR EACH ROW EXECUTE PROCEDURE nft_update_updated_at_column();

CREATE TABLE IF NOT EXISTS nft_categories (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR NOT NULL,
    sort BIGINT NOT NULL DEFAULT 0,
    main boolean DEFAULT false
);
CREATE INDEX ordering_nft_categories_list ON nft_categories USING BTREE (sort);

CREATE TABLE IF NOT EXISTS nft_relations (
    nft_item_id uuid NOT NULL,
    relation_id uuid NOT NULL,
    PRIMARY KEY(nft_item_id, relation_id),
    FOREIGN KEY(nft_item_id) REFERENCES nft_items(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TRIGGER IF EXISTS update_nft_items_modtime ON nft_items;
DROP TABLE IF EXISTS nft_items;
DROP TABLE IF EXISTS nft_categories;
DROP TABLE IF EXISTS nft_relations;
DROP FUNCTION IF EXISTS nft_update_updated_at_column();