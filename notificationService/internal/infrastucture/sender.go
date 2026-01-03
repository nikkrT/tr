package infrastucture

import (
	"context"
	"encoding/json"
	"fmt"
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

type Sender struct {
}

func NewInfrastucture() *Sender {
	return &Sender{}
}

func (s *Sender) EmailCreated(ctx context.Context, data []byte) error {
	product := Product{}
	err := json.Unmarshal(data, &product)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal the data: %s", err)
	}
	fmt.Println("product created", product)
	return nil
}

func (s *Sender) EmailUpdated(ctx context.Context, data []byte) error {
	product := Product{}
	err := json.Unmarshal(data, &product)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal the data: %s", err)
	}
	fmt.Println("product updated", product)
	return nil
}
func (s *Sender) EmailDeleted(ctx context.Context, data []byte) error {
	var id int
	err := json.Unmarshal(data, &id)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal the data: %s", err)
	}
	fmt.Println("product deleted with id", id)
	return nil
}
