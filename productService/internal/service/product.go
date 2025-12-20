package service

import (
	"context"
	"errors"
	"fmt"
	"productService/internal/model"

	"github.com/go-playground/validator/v10"
)

type ProductRepository interface {
	Create(ctx context.Context, p model.Product) error
	ReadById(ctx context.Context, id int) (model.Product, error)
	ReadAll(ctx context.Context, filtered string) ([]model.Product, error)
	Update(ctx context.Context, p model.Product) error
	DeleteById(ctx context.Context, id int) error
}
type ProductService struct {
	Repo ProductRepository
}

func NewProductService(repository ProductRepository) *ProductService {
	return &ProductService{Repo: repository}
}

func ValidateProduct(product model.Product) error {
	validate := validator.New()
	return validate.Struct(product)
}

func (s *ProductService) CreateProduct(ctx context.Context, p model.Product) error {
	err := ValidateProduct(p)
	if err != nil {
		return fmt.Errorf("Product validation error: %w", err)
	}
	return s.Repo.Create(ctx, p)
}

var ErrNotFound = errors.New("product not found")

func (s *ProductService) ReadProduct(ctx context.Context, id int) (model.Product, error) {
	product, err := s.Repo.ReadById(ctx, id)
	if err != nil {
		return model.Product{}, fmt.Errorf("Product read error: %w", ErrNotFound)
	}
	return product, nil
}
func (s *ProductService) ReadAll(ctx context.Context, filteredBy string) ([]model.Product, error) {
	products, err := s.Repo.ReadAll(ctx, filteredBy)
	if err != nil {
		return []model.Product{}, fmt.Errorf("Product read error: %w", err)
	}
	return products, nil
}
func (s *ProductService) UpdateProduct(ctx context.Context, p model.Product) error {
	err := ValidateProduct(p)
	if err != nil {
		return fmt.Errorf("Product validation error: %w", err)
	}
	return s.Repo.Update(ctx, p)
}
func (s *ProductService) DeleteProduct(ctx context.Context, id int) error {
	err := s.Repo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("Product delete error: %w", err)
	}
	return nil
}
