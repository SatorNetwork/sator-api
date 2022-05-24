// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"database/sql"
	"time"
)

type SolanaError struct {
	ProviderName string       `json:"provider_name"`
	ErrorMessage string       `json:"error_message"`
	Counter      int32        `json:"counter"`
	UpdatedAt    sql.NullTime `json:"updated_at"`
	CreatedAt    time.Time    `json:"created_at"`
}

type SolanaMetric struct {
	ProviderName       string       `json:"provider_name"`
	NotAvailableErrors int32        `json:"not_available_errors"`
	OtherErrors        int32        `json:"other_errors"`
	SuccessCalls       int32        `json:"success_calls"`
	UpdatedAt          sql.NullTime `json:"updated_at"`
	CreatedAt          time.Time    `json:"created_at"`
}
