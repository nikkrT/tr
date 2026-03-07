package repo

import (
	"context"
	"fmt"
	"micr_course/pkg/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo struct {
	psql *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{psql: db}
}

func (h *OrderRepo) Create(ctx context.Context, order models.OrderModel) (int, error) {
	sql := "INSERT INTO orders (product, status, created) VALUES ($1, $2, $3) RETURNING id"
	var id int
	err := h.psql.QueryRow(ctx, sql, order.ProductId, order.Status, time.Now()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting product into database: %w", err)
	}
	return id, nil
}

func (h *OrderRepo) ReadById(ctx context.Context, id int) (models.OrderModel, error) {
	sql := "SELECT * FROM orders WHERE id=$1"
	order := models.OrderModel{}
	err := h.psql.QueryRow(ctx, sql, id).Scan(&order.Id, &order.ProductId, &order.Status, &order.CreatedAt, &order.DeletedAt)
	if err != nil {
		return models.OrderModel{}, fmt.Errorf("scan err: %w", err)
	}
	if order.Id == 0 {
		return models.OrderModel{}, fmt.Errorf("Product not found")
	}
	return order, nil
}

func (h *OrderRepo) Update(ctx context.Context, order models.OrderModel) error {

	sql := "UPDATE orders " +
		"SET status=$1 WHERE id=$2"
	rows, err := h.psql.Exec(ctx, sql, &order.Status, &order.Id)
	if err != nil {
		return fmt.Errorf("update product rr: %w", err)
	}
	affected := rows.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("update product rr: %w", pgx.ErrNoRows)
	}
	return nil
}

func (h *OrderRepo) DeleteById(ctx context.Context, id int) error {
	sql := "DELETE FROM order WHERE id=$1"
	rows, err := h.psql.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("delete order rr: %w", err)
	}
	affected := rows.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("delete order rr: %w", pgx.ErrNoRows)
	}
	return nil
}
