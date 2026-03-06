package notificator

import (
	"context"
	"encoding/json"
	"fmt"
	model "micr_course/pkg/models"
)

type Sender struct {
}

func NewInfrastucture() *Sender {
	return &Sender{}
}

func (s *Sender) EmailCreated(ctx context.Context, data []byte) error {
	product := model.Product{}
	err := json.Unmarshal(data, &product)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal the data: %s", err)
	}
	fmt.Println("product created", product)
	return nil
}

func (s *Sender) EmailUpdated(ctx context.Context, data []byte) error {
	product := model.Product{}
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
