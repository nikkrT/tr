package postgres

import (
	"context"
	"fmt"
	"micr_course/productService/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	//err = createTables(ctx, pool)
	//if err != nil {
	//	return nil, fmt.Errorf("error creating tables: %v", err)
	//}
	return pool, nil
}
