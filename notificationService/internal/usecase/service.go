package usecase

import (
	"context"
	"fmt"
)

type Infrastucture interface {
	EmailCreated(ctx context.Context, data []byte) error
	EmailUpdated(ctx context.Context, data []byte) error
	EmailDeleted(ctx context.Context, data []byte) error
}

type Service struct {
	Infrastucture Infrastucture
}

func NewService(infrastucture Infrastucture) *Service {
	return &Service{Infrastucture: infrastucture}
}

func (s *Service) SendCreated(ctx context.Context, data []byte) error {
	err := s.Infrastucture.EmailCreated(ctx, data)
	if err != nil {
		return fmt.Errorf("send email created failed: %w", err)
	}
	return nil
}
func (s *Service) SendUpdated(ctx context.Context, data []byte) error {
	err := s.Infrastucture.EmailUpdated(ctx, data)
	if err != nil {
		return fmt.Errorf("send email updated failed: %w", err)
	}
	return nil
}
func (s *Service) SendDeleted(ctx context.Context, data []byte) error {
	err := s.Infrastucture.EmailDeleted(ctx, data)
	if err != nil {
		return fmt.Errorf("send email deleted failed: %w", err)
	}
	return nil
}
