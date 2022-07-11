// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UnityGameSettingsValueType string

const (
	UnityGameSettingsValueTypeString   UnityGameSettingsValueType = "string"
	UnityGameSettingsValueTypeInt      UnityGameSettingsValueType = "int"
	UnityGameSettingsValueTypeFloat    UnityGameSettingsValueType = "float"
	UnityGameSettingsValueTypeBool     UnityGameSettingsValueType = "bool"
	UnityGameSettingsValueTypeJson     UnityGameSettingsValueType = "json"
	UnityGameSettingsValueTypeDuration UnityGameSettingsValueType = "duration"
	UnityGameSettingsValueTypeDatetime UnityGameSettingsValueType = "datetime"
)

func (e *UnityGameSettingsValueType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UnityGameSettingsValueType(s)
	case string:
		*e = UnityGameSettingsValueType(s)
	default:
		return fmt.Errorf("unsupported scan type for UnityGameSettingsValueType: %T", src)
	}
	return nil
}

type UnityGameNft struct {
	ID           string         `json:"id"`
	UserID       uuid.UUID      `json:"user_id"`
	NftType      string         `json:"nft_type"`
	MaxLevel     int32          `json:"max_level"`
	CraftedNftID sql.NullString `json:"crafted_nft_id"`
	DeletedAt    sql.NullTime   `json:"deleted_at"`
}

type UnityGameNftPack struct {
	ID          uuid.UUID    `json:"id"`
	DropChances []byte       `json:"drop_chances"`
	Price       float64      `json:"price"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
	CreatedAt   time.Time    `json:"created_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
	Name        string       `json:"name"`
}

type UnityGamePlayer struct {
	UserID              uuid.UUID      `json:"user_id"`
	EnergyPoints        int32          `json:"energy_points"`
	EnergyRefilledAt    time.Time      `json:"energy_refilled_at"`
	SelectedNftID       sql.NullString `json:"selected_nft_id"`
	UpdatedAt           sql.NullTime   `json:"updated_at"`
	CreatedAt           time.Time      `json:"created_at"`
	ElectricitySpent    int32          `json:"electricity_spent"`
	ElectricityCosts    float64        `json:"electricity_costs"`
	SelectedCharacterID sql.NullString `json:"selected_character_id"`
}

type UnityGameResult struct {
	ID               uuid.UUID     `json:"id"`
	UserID           uuid.UUID     `json:"user_id"`
	NFTID            string        `json:"nft_id"`
	Complexity       int32         `json:"complexity"`
	IsTraining       bool          `json:"is_training"`
	BlocksDone       int32         `json:"blocks_done"`
	FinishedAt       sql.NullTime  `json:"finished_at"`
	UpdatedAt        sql.NullTime  `json:"updated_at"`
	CreatedAt        time.Time     `json:"created_at"`
	Result           sql.NullInt32 `json:"result"`
	ElectricityCosts float64       `json:"electricity_costs"`
}

type UnityGameReward struct {
	ID            uuid.UUID     `json:"id"`
	UserID        uuid.UUID     `json:"user_id"`
	RelationID    uuid.NullUUID `json:"relation_id"`
	OperationType int32         `json:"operation_type"`
	Amount        float64       `json:"amount"`
	CreatedAt     time.Time     `json:"created_at"`
}

type UnityGameSetting struct {
	Key         string                     `json:"key"`
	Name        string                     `json:"name"`
	ValueType   UnityGameSettingsValueType `json:"value_type"`
	Value       string                     `json:"value"`
	Description sql.NullString             `json:"description"`
}
