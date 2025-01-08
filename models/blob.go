package models

import (
	"time"
	"github.com/google/uuid"
)

type Blob struct {
	BlobID      uuid.UUID `json:"blob_id" db:"blob_id" validate:"omitempty,uuid"`
	UserID      string    `json:"user_id" db:"user_id" validate:"required"`
	Description string    `json:"description" db:"description" validate:"required,lte=15"`
	Title       string    `json:"title" db:"title" validate:"required,lte=10"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type BlobList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*Blob `json:"blobs"`
}
