// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: nfts.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const addNFT = `-- name: AddNFT :one
INSERT INTO unity_game_nfts (id, user_id, nft_type, max_level)  
VALUES ($1, $2, $3, $4) RETURNING id, user_id, nft_type, max_level, crafted_nft_id, deleted_at
`

type AddNFTParams struct {
	ID       string    `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	NftType  string    `json:"nft_type"`
	MaxLevel int32     `json:"max_level"`
}

func (q *Queries) AddNFT(ctx context.Context, arg AddNFTParams) (UnityGameNft, error) {
	row := q.queryRow(ctx, q.addNFTStmt, addNFT,
		arg.ID,
		arg.UserID,
		arg.NftType,
		arg.MaxLevel,
	)
	var i UnityGameNft
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.NftType,
		&i.MaxLevel,
		&i.CraftedNftID,
		&i.DeletedAt,
	)
	return i, err
}

const craftNFTs = `-- name: CraftNFTs :exec
UPDATE unity_game_nfts SET crafted_nft_id = $1, deleted_at = now()
WHERE user_id = $2 AND id = ANY($3::VARCHAR[])
`

type CraftNFTsParams struct {
	CraftedNftID sql.NullString `json:"crafted_nft_id"`
	UserID       uuid.UUID      `json:"user_id"`
	NftIds       []string       `json:"nft_ids"`
}

func (q *Queries) CraftNFTs(ctx context.Context, arg CraftNFTsParams) error {
	_, err := q.exec(ctx, q.craftNFTsStmt, craftNFTs, arg.CraftedNftID, arg.UserID, pq.Array(arg.NftIds))
	return err
}

const deleteNFT = `-- name: DeleteNFT :exec
UPDATE unity_game_nfts SET deleted_at = now() WHERE id = $1
`

func (q *Queries) DeleteNFT(ctx context.Context, id string) error {
	_, err := q.exec(ctx, q.deleteNFTStmt, deleteNFT, id)
	return err
}

const getNFT = `-- name: GetNFT :one
SELECT id, user_id, nft_type, max_level, crafted_nft_id, deleted_at FROM unity_game_nfts WHERE id = $1
`

func (q *Queries) GetNFT(ctx context.Context, id string) (UnityGameNft, error) {
	row := q.queryRow(ctx, q.getNFTStmt, getNFT, id)
	var i UnityGameNft
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.NftType,
		&i.MaxLevel,
		&i.CraftedNftID,
		&i.DeletedAt,
	)
	return i, err
}

const getUserNFT = `-- name: GetUserNFT :one
SELECT id, user_id, nft_type, max_level, crafted_nft_id, deleted_at FROM unity_game_nfts WHERE id = $1 AND user_id = $2
`

type GetUserNFTParams struct {
	ID     string    `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) GetUserNFT(ctx context.Context, arg GetUserNFTParams) (UnityGameNft, error) {
	row := q.queryRow(ctx, q.getUserNFTStmt, getUserNFT, arg.ID, arg.UserID)
	var i UnityGameNft
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.NftType,
		&i.MaxLevel,
		&i.CraftedNftID,
		&i.DeletedAt,
	)
	return i, err
}

const getUserNFTByIDs = `-- name: GetUserNFTByIDs :many
SELECT id, user_id, nft_type, max_level, crafted_nft_id, deleted_at FROM unity_game_nfts 
WHERE user_id = $1
AND id = ANY($2::VARCHAR[])
AND deleted_at IS NULL 
AND crafted_nft_id IS NULL
`

type GetUserNFTByIDsParams struct {
	UserID uuid.UUID `json:"user_id"`
	IDs    []string  `json:"ids"`
}

func (q *Queries) GetUserNFTByIDs(ctx context.Context, arg GetUserNFTByIDsParams) ([]UnityGameNft, error) {
	rows, err := q.query(ctx, q.getUserNFTByIDsStmt, getUserNFTByIDs, arg.UserID, pq.Array(arg.IDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UnityGameNft
	for rows.Next() {
		var i UnityGameNft
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.NftType,
			&i.MaxLevel,
			&i.CraftedNftID,
			&i.DeletedAt,
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

const getUserNFTs = `-- name: GetUserNFTs :many
SELECT id, user_id, nft_type, max_level, crafted_nft_id, deleted_at FROM unity_game_nfts 
WHERE user_id = $1 
AND deleted_at IS NULL 
AND crafted_nft_id IS NULL
`

func (q *Queries) GetUserNFTs(ctx context.Context, userID uuid.UUID) ([]UnityGameNft, error) {
	rows, err := q.query(ctx, q.getUserNFTsStmt, getUserNFTs, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UnityGameNft
	for rows.Next() {
		var i UnityGameNft
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.NftType,
			&i.MaxLevel,
			&i.CraftedNftID,
			&i.DeletedAt,
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
