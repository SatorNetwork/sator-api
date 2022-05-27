// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: players.sql

package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const addElectricityToPlayer = `-- name: AddElectricityToPlayer :exec
UPDATE unity_game_players 
SET 
    electricity_costs = electricity_costs + $1,
    electricity_spent = electricity_spent + 1 
WHERE user_id = $2
`

type AddElectricityToPlayerParams struct {
	ElectricityCosts float64   `json:"electricity_costs"`
	UserID           uuid.UUID `json:"user_id"`
}

func (q *Queries) AddElectricityToPlayer(ctx context.Context, arg AddElectricityToPlayerParams) error {
	_, err := q.exec(ctx, q.addElectricityToPlayerStmt, addElectricityToPlayer, arg.ElectricityCosts, arg.UserID)
	return err
}

const addNewPlayer = `-- name: AddNewPlayer :one
INSERT INTO unity_game_players (user_id, energy_points, energy_refilled_at, selected_nft_id) 
VALUES ($1, $2, $3, $4) RETURNING user_id, energy_points, energy_refilled_at, selected_nft_id, updated_at, created_at, electricity_spent, electricity_costs
`

type AddNewPlayerParams struct {
	UserID           uuid.UUID      `json:"user_id"`
	EnergyPoints     int32          `json:"energy_points"`
	EnergyRefilledAt time.Time      `json:"energy_refilled_at"`
	SelectedNftID    sql.NullString `json:"selected_nft_id"`
}

func (q *Queries) AddNewPlayer(ctx context.Context, arg AddNewPlayerParams) (UnityGamePlayer, error) {
	row := q.queryRow(ctx, q.addNewPlayerStmt, addNewPlayer,
		arg.UserID,
		arg.EnergyPoints,
		arg.EnergyRefilledAt,
		arg.SelectedNftID,
	)
	var i UnityGamePlayer
	err := row.Scan(
		&i.UserID,
		&i.EnergyPoints,
		&i.EnergyRefilledAt,
		&i.SelectedNftID,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.ElectricitySpent,
		&i.ElectricityCosts,
	)
	return i, err
}

const getPlayer = `-- name: GetPlayer :one
SELECT user_id, energy_points, energy_refilled_at, selected_nft_id, updated_at, created_at, electricity_spent, electricity_costs FROM unity_game_players WHERE user_id = $1
`

func (q *Queries) GetPlayer(ctx context.Context, userID uuid.UUID) (UnityGamePlayer, error) {
	row := q.queryRow(ctx, q.getPlayerStmt, getPlayer, userID)
	var i UnityGamePlayer
	err := row.Scan(
		&i.UserID,
		&i.EnergyPoints,
		&i.EnergyRefilledAt,
		&i.SelectedNftID,
		&i.UpdatedAt,
		&i.CreatedAt,
		&i.ElectricitySpent,
		&i.ElectricityCosts,
	)
	return i, err
}

const refillEnergyOfPlayer = `-- name: RefillEnergyOfPlayer :exec
UPDATE unity_game_players SET energy_points = $1, energy_refilled_at = $2 WHERE user_id = $3
`

type RefillEnergyOfPlayerParams struct {
	EnergyPoints     int32     `json:"energy_points"`
	EnergyRefilledAt time.Time `json:"energy_refilled_at"`
	UserID           uuid.UUID `json:"user_id"`
}

func (q *Queries) RefillEnergyOfPlayer(ctx context.Context, arg RefillEnergyOfPlayerParams) error {
	_, err := q.exec(ctx, q.refillEnergyOfPlayerStmt, refillEnergyOfPlayer, arg.EnergyPoints, arg.EnergyRefilledAt, arg.UserID)
	return err
}

const resetElectricityForPlayer = `-- name: ResetElectricityForPlayer :exec
UPDATE unity_game_players SET electricity_costs = 0, electricity_spent = 0 WHERE user_id = $1
`

func (q *Queries) ResetElectricityForPlayer(ctx context.Context, userID uuid.UUID) error {
	_, err := q.exec(ctx, q.resetElectricityForPlayerStmt, resetElectricityForPlayer, userID)
	return err
}

const spendEnergyOfPlayer = `-- name: SpendEnergyOfPlayer :exec
UPDATE unity_game_players SET energy_points = energy_points-1 WHERE user_id = $1
`

func (q *Queries) SpendEnergyOfPlayer(ctx context.Context, userID uuid.UUID) error {
	_, err := q.exec(ctx, q.spendEnergyOfPlayerStmt, spendEnergyOfPlayer, userID)
	return err
}

const storeSelectedNFT = `-- name: StoreSelectedNFT :exec
UPDATE unity_game_players SET selected_nft_id = $1 WHERE user_id = $2
`

type StoreSelectedNFTParams struct {
	SelectedNftID sql.NullString `json:"selected_nft_id"`
	UserID        uuid.UUID      `json:"user_id"`
}

func (q *Queries) StoreSelectedNFT(ctx context.Context, arg StoreSelectedNFTParams) error {
	_, err := q.exec(ctx, q.storeSelectedNFTStmt, storeSelectedNFT, arg.SelectedNftID, arg.UserID)
	return err
}
