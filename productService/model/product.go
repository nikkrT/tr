package model

import (
	"time"
)

type Product struct {
	Id          int        `json:"id"`
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Price       int        `json:"price" validate:"required"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
