package models

import (
	"time"

	"github.com/google/uuid"
)

type Interest struct {
	ID        	uuid.UUID     `json:"id" db:"id" validate:"required,uuid"`
	Name        string    `json:"name" db:"name" validate:"required"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
