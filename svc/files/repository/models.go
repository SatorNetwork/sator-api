// Code generated by sqlc. DO NOT EDIT.

package repository

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        uuid.UUID `json:"id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	FileUrl   string    `json:"file_url"`
	CreatedAt time.Time `json:"created_at"`
}