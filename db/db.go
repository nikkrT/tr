package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

//func createTables(ctx context.Context, db *pgxpool.Pool) error {
//	tables := []struct {
//		name string
//		sql  string
//	}{
//		{
//			name: "products",
//			sql: `CREATE TABLE IF NOT EXISTS products (
//					id SERIAL PRIMARY KEY,
//					name VARCHAR(255) NOT NULL,
//					description TEXT,
//					price INT NOT NULL,
//					created_at TIMESTAMPTZ DEFAULT NOW(),
//					updated_at TIMESTAMPTZ DEFAULT NOW(),
//					deleted_at TIMESTAMPTZ
//            );`,
//		},
//	}
//	for _, table := range tables {
//		if _, err := db.Exec(ctx, table.sql); err != nil {
//			return fmt.Errorf("create table %s: %w", table.name, err)
//		}
//	}
//	return nil
//}

func InitDB(ctx context.Context, URL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	defer pool.Close()
	//err = createTables(ctx, pool)
	//if err != nil {
	//	return nil, fmt.Errorf("error creating tables: %v", err)
	//}
	return pool, nil
}
