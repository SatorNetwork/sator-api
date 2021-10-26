// Code generated by sqlc. DO NOT EDIT.
// source: items.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const addNFTItem = `-- name: AddNFTItem :one
INSERT INTO nft_items (name, description, cover, supply, buy_now_price, token_uri)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, owner_id, name, description, cover, supply, buy_now_price, token_uri, updated_at, created_at
`

type AddNFTItemParams struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cover       string  `json:"cover"`
	Supply      int64   `json:"supply"`
	BuyNowPrice float64 `json:"buy_now_price"`
	TokenURI    string  `json:"token_uri"`
}

func (q *Queries) AddNFTItem(ctx context.Context, arg AddNFTItemParams) (NFTItem, error) {
	row := q.queryRow(ctx, q.addNFTItemStmt, addNFTItem,
		arg.Name,
		arg.Description,
		arg.Cover,
		arg.Supply,
		arg.BuyNowPrice,
		arg.TokenURI,
	)
	var i NFTItem
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Name,
		&i.Description,
		&i.Cover,
		&i.Supply,
		&i.BuyNowPrice,
		&i.TokenURI,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getNFTItemByID = `-- name: GetNFTItemByID :one
SELECT id, owner_id, name, description, cover, supply, buy_now_price, token_uri, updated_at, created_at FROM nft_items
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetNFTItemByID(ctx context.Context, id uuid.UUID) (NFTItem, error) {
	row := q.queryRow(ctx, q.getNFTItemByIDStmt, getNFTItemByID, id)
	var i NFTItem
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Name,
		&i.Description,
		&i.Cover,
		&i.Supply,
		&i.BuyNowPrice,
		&i.TokenURI,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getNFTItemsList = `-- name: GetNFTItemsList :many
SELECT id, owner_id, name, description, cover, supply, buy_now_price, token_uri, updated_at, created_at FROM nft_items
ORDER BY updated_at DESC, created_at DESC
LIMIT $2 OFFSET $1
`

type GetNFTItemsListParams struct {
	Offset int32 `json:"offset_val"`
	Limit  int32 `json:"limit_val"`
}

func (q *Queries) GetNFTItemsList(ctx context.Context, arg GetNFTItemsListParams) ([]NFTItem, error) {
	rows, err := q.query(ctx, q.getNFTItemsListStmt, getNFTItemsList, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []NFTItem
	for rows.Next() {
		var i NFTItem
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.Name,
			&i.Description,
			&i.Cover,
			&i.Supply,
			&i.BuyNowPrice,
			&i.TokenURI,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNFTItemsListByOwnerID = `-- name: GetNFTItemsListByOwnerID :many
SELECT id, owner_id, name, description, cover, supply, buy_now_price, token_uri, updated_at, created_at FROM nft_items
WHERE owner_id = $1
ORDER BY updated_at DESC, created_at DESC
LIMIT $3 OFFSET $2
`

type GetNFTItemsListByOwnerIDParams struct {
	OwnerID uuid.NullUUID `json:"owner_id"`
	Offset  int32         `json:"offset_val"`
	Limit   int32         `json:"limit_val"`
}

func (q *Queries) GetNFTItemsListByOwnerID(ctx context.Context, arg GetNFTItemsListByOwnerIDParams) ([]NFTItem, error) {
	rows, err := q.query(ctx, q.getNFTItemsListByOwnerIDStmt, getNFTItemsListByOwnerID, arg.OwnerID, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []NFTItem
	for rows.Next() {
		var i NFTItem
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.Name,
			&i.Description,
			&i.Cover,
			&i.Supply,
			&i.BuyNowPrice,
			&i.TokenURI,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNFTItemsListByRelationID = `-- name: GetNFTItemsListByRelationID :many
SELECT nft_items.id, nft_items.owner_id, nft_items.name, nft_items.description, nft_items.cover, nft_items.supply, nft_items.buy_now_price, nft_items.token_uri, nft_items.updated_at, nft_items.created_at FROM nft_items
JOIN nft_relations ON nft_relations.nft_item_id =  nft_items.id
WHERE nft_relations.relation_id = $1
AND owner_id = $2
ORDER BY created_at DESC
LIMIT $4 OFFSET $3
`

type GetNFTItemsListByRelationIDParams struct {
	RelationID uuid.UUID     `json:"relation_id"`
	OwnerID    uuid.NullUUID `json:"owner_id"`
	Offset     int32         `json:"offset_val"`
	Limit      int32         `json:"limit_val"`
}

func (q *Queries) GetNFTItemsListByRelationID(ctx context.Context, arg GetNFTItemsListByRelationIDParams) ([]NFTItem, error) {
	rows, err := q.query(ctx, q.getNFTItemsListByRelationIDStmt, getNFTItemsListByRelationID,
		arg.RelationID,
		arg.OwnerID,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []NFTItem
	for rows.Next() {
		var i NFTItem
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.Name,
			&i.Description,
			&i.Cover,
			&i.Supply,
			&i.BuyNowPrice,
			&i.TokenURI,
			&i.UpdatedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateNFTItemOwner = `-- name: UpdateNFTItemOwner :exec
UPDATE nft_items SET owner_id = $1
WHERE id = $2
`

type UpdateNFTItemOwnerParams struct {
	OwnerID uuid.NullUUID `json:"owner_id"`
	ID      uuid.UUID     `json:"id"`
}

func (q *Queries) UpdateNFTItemOwner(ctx context.Context, arg UpdateNFTItemOwnerParams) error {
	_, err := q.exec(ctx, q.updateNFTItemOwnerStmt, updateNFTItemOwner, arg.OwnerID, arg.ID)
	return err
}
