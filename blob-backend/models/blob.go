package models

import (
	"time"
	"github.com/google/uuid"
)

type Blob struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"omitempty"`
	Content   string    `json:"content" db:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserID    string    `json:"user_id" db:"user_id" validate:"required,uuid"`
}

type BlobWithInterests struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"omitempty"`
	Content   string    `json:"content" db:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UserID    string    `json:"user_id" db:"user_id" validate:"required,uuid"`
	Interests     []string  `json:"interests"`
}

type BlobListWithDetails struct {
    ID            string 	`json:"id" db:"id"`
    UserID        string    `json:"user_id" db:"user_id"`
    Content       string    `json:"content" db:"content"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
    Username      string    `json:"username" db:"username"`
    AvatarIcon    string    `json:"avatar_icon" db:"avatar_icon"`
    UserCreatedAt time.Time `json:"user_created_at" db:"user_created_at"`
    LikesCount    int       `json:"likes_count" db:"likes_count"`
    CommentsCount int       `json:"comments_count" db:"comments_count"`
    Interests     []string  `json:"interests"`
}

type BlobWithDetails struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Username     string    `json:"username"`
	AvatarIcon   string    `json:"avatar_icon"`
	UserCreatedAt time.Time `json:"user_created_at"`
	Comments     []Comment `json:"comments"`
	Likes        []Like    `json:"likes"`
	Interests    []Interest  `json:"interests"`
}

type BlobList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*BlobListWithDetails `json:"blobs"`
}
