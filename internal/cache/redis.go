package cache

import (
	"context"
	"encoding/json"
	"wb-tech/internal/common/config"
	"wb-tech/internal/model"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedis(cfg config.Redis) *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Host + ":" + cfg.Port,
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
	}
}

func (r *RedisClient) GetOrder(ctx context.Context, uid string) (*model.Order, error) {
	data, err := r.client.Get(ctx, uid).Bytes()
	if err != nil {
		return nil, err
	}

	var order model.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *RedisClient) SetOrder(ctx context.Context, order *model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, order.OrderUID, data, 0).Err()
}
