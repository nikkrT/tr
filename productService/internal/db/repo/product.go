package repo

import (
	"context"
	"fmt"
	"productService/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	DB *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (h *ProductRepo) Create(ctx context.Context, product model.Product) error {
	sql := "INSERT INTO products (name,description, price,created_at) VALUES ($1, $2, $3, $4)"
	_, err := h.DB.Exec(ctx, sql, product.Name, product.Description, product.Price, time.Now())
	if err != nil {
		return fmt.Errorf("Error inserting product into database: %w", err)
	}
	return nil
}

func (h *ProductRepo) ReadById(ctx context.Context, id int) (model.Product, error) {
	sql := "SELECT * FROM products WHERE id=$1"
	product := model.Product{}
	err := h.DB.QueryRow(ctx, sql, id).Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)
	if err != nil {
		return model.Product{}, fmt.Errorf("scan err: %w", err)
	}
	if product.Id == 0 {
		return model.Product{}, fmt.Errorf("Product not found")
	}
	return product, nil
}

func (h *ProductRepo) ReadAll(ctx context.Context, filtered string) ([]model.Product, error) {

	allowedFilters := map[string]string{
		"id":         "id",
		"name":       "name",
		"price":      "price",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}
	column, ok := allowedFilters[filtered]
	if !ok {
		column = "id"
	}
	sql := fmt.Sprintf("SELECT * FROM products ORDER BY %s", column)
	rows, err := h.DB.Query(ctx, sql)
	fmt.Println(rows)
	if err != nil {
		return nil, fmt.Errorf("find productSS rr: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		product := model.Product{}
		err = rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("scan err: %w", err)
		}
		products = append(products, product)
	}
	fmt.Println(products)
	return products, nil
}

func (h *ProductRepo) Update(ctx context.Context, product model.Product) error {

	sql := "UPDATE products " +
		"SET description=$1, price=$2, updated_at=$3 WHERE id=$4"
	rows, err := h.DB.Exec(ctx, sql, &product.Description, &product.Price, time.Now(), &product.Id)
	if err != nil {
		return fmt.Errorf("update product rr: %w", err)
	}
	affected := rows.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("update product rr: %w", pgx.ErrNoRows)
	}
	return nil
}

func (h *ProductRepo) DeleteById(ctx context.Context, id int) error {
	sql := "DELETE FROM products WHERE id=$1"
	rows, err := h.DB.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete product rr: %w", err)
	}
	affected := rows.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("delete product rr: %w", pgx.ErrNoRows)
	}
	return nil
}
