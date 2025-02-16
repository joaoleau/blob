package models

import (
	"time"
	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID     `json:"id" db:"id" validate:"required,uuid"`
	Content   string        `json:"content" db:"content" validate:"required"`
	CreatedAt time.Time   	`json:"created_at" db:"created_at"`
	UpdatedAt time.Time   	`json:"updated_at" db:"updated_at"`
	UserID    string       	`json:"user_id" db:"user_id" validate:"required,uuid"`
	BlobID    uuid.UUID    	`json:"blob_id" db:"blob_id" validate:"required,uuid"`
}

type CommentWithUser struct {
	ID        uuid.UUID     `json:"id" db:"id" validate:"required,uuid"`
	Content   string        `json:"content" db:"content" validate:"required"`
	CreatedAt time.Time   	`json:"created_at" db:"created_at"`
	UpdatedAt time.Time   	`json:"updated_at" db:"updated_at"`
	UserID    		string  `json:"user_id" db:"user_id" validate:"required,uuid"`
	Image         string    `json:"image,omitempty" db:"image"`
	Username      string    `json:"username,omitempty" db:"username"`
	AvatarIcon    string    `json:"avatar_icon" db:"avatar_icon" default:"user"`
	AvatarColor   string    `json:"avatar_color" db:"avatar_color" default:"cyan"`
	BlobID    uuid.UUID    	`json:"blob_id" db:"blob_id" validate:"required,uuid"`
}
