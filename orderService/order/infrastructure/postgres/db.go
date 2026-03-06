package db

import (
	"context"
	"fmt"
	"micr_course/orderService/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(ctx context.Context, cfg config.Postgres) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.Address)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return pool, nil
}
