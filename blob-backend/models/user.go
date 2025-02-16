package models

import (
	"time"
)

type User struct {
	ID            string    `json:"id" db:"id" validate:"required,uuid"`
	Name          string    `json:"name,omitempty" db:"name"`
	Email         string    `json:"email,omitempty" db:"email" validate:"omitempty,email"`
	EmailVerified time.Time `json:"email_verified,omitempty" db:"email_verified"`
	Image         string    `json:"image,omitempty" db:"image"`
	Password      string    `json:"password,omitempty" db:"password"`
	Username      string    `json:"username,omitempty" db:"username"`
	Bio           string    `json:"bio,omitempty" db:"bio"`
	AvatarIcon    string    `json:"avatar_icon" db:"avatar_icon" default:"user"`
	AvatarColor   string    `json:"avatar_color" db:"avatar_color" default:"cyan"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type UserWithBlobs struct {
	ID              string    `db:"id"`
	Name            string    `db:"name"`
	Email           string    `db:"email"`
	EmailVerified   time.Time `db:"email_verified"`
	Image           string    `db:"image"`
	Username        string    `db:"username"`
	Bio             string    `db:"bio"`
	AvatarIcon      string    `db:"avatar_icon"`
	AvatarColor     string    `db:"avatar_color"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	Blobs           []Blob    `db:"-"`
}

type UserList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*User `json:"users"`
}