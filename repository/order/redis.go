package order

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"tr/model"
)

type RedisRepository struct {
	Client *redis.Client
}

func OrderId(id uint64) string {
	return fmt.Sprintf("order:%v", id)
}

func (r *RedisRepository) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	key := OrderId(order.OrderId)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, data, time.Hour)

	if res.Err() != nil {
		txn.Discard()
		return res.Err()
	}

	if err := txn.SAdd(ctx, "order", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("не получилось в сет добавить", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}

var notExistErr = errors.New("order not exist")

func (r *RedisRepository) FindById(ctx context.Context, id uint64) (*model.Order, error) {
	key := OrderId(id)
	res, err := r.Client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, notExistErr
	}
	order := &model.Order{}
	err = json.Unmarshal([]byte(res), &order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *RedisRepository) DeleteById(ctx context.Context, id uint64) error {
	key := OrderId(id)
	res := r.Client.Del(ctx, key)
	if errors.Is(res.Err(), redis.Nil) {
		return notExistErr
	}
	return nil
}

func (r *RedisRepository) UpdateById(ctx context.Context, order model.Order) error {
	key := OrderId(order.OrderId)
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}
	err = r.Client.SetXX(ctx, key, string(data), time.Hour).Err()
	if errors.Is(err, redis.Nil) {
		return notExistErr
	} else if err != nil {
		return fmt.Errorf("failed to set order: %w", err)
	}
	return nil
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepository) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, "orders", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get order ids: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := make([]model.Order, len(xs))

	for i, x := range xs {
		x := x.(string)
		var order model.Order

		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode order json: %w", err)
		}

		orders[i] = order
	}

	return FindResult{
		Orders: orders,
		Cursor: cursor,
	}, nil
}
