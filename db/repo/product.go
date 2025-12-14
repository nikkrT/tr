package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"micr_course/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	DB *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (w *ProductRepo) FindById(ctx context.Context, id int) (model.Product, error) {
	sql := "SELECT * FROM products WHERE id=$1"
	rows, err := w.DB.Query(ctx, sql, id)
	defer rows.Close()
	if err != nil {
		return model.Product{}, fmt.Errorf("find product rr: %w", err)
	}
	product := model.Product{}
	if rows.Next() {
		err = rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)
		if err != nil {
			return model.Product{}, fmt.Errorf("scan product rr: %w", err)
		}
		data, _ := json.Marshal(product)
		fmt.Println(string(data))
		fmt.Println(product)
		return product, nil
	}
	return model.Product{}, fmt.Errorf("find product rr: %w", pgx.ErrNoRows)
}
