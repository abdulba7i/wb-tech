package cache

import (
	"context"
	"errors"
	"sync"
	"wb-tech/internal/model"
	"wb-tech/internal/repository"
)

var (
	ErrNotFound = errors.New("order not found")
)

type OrderCache interface {
	Get(ctx context.Context, uid string) (*model.Order, error)
	Set(ctx context.Context, order *model.Order) error
	Preload(ctx context.Context) error
}

type CacheService struct {
	repo        repository.OrderRepository
	Redis       *RedisClient
	preloadOnce sync.Once
}

func NewCacheService(repo repository.OrderRepository, redis *RedisClient) *CacheService {
	return &CacheService{
		repo:  repo,
		Redis: redis,
	}
}

func (s *CacheService) Get(ctx context.Context, uid string) (*model.Order, error) {
	order, err := s.Redis.GetOrder(ctx, uid)
	if err == nil && order != nil {
		return order, nil
	}

	order, err = s.repo.GetOrderByUID(uid)
	if err != nil {
		return nil, err
	}
	s.Redis.SetOrder(ctx, order)
	return order, nil
}

func (s *CacheService) Set(ctx context.Context, order *model.Order) error {
	return s.Redis.SetOrder(ctx, order)
}

func (s *CacheService) Preload(ctx context.Context) error {
	var initErr error
	s.preloadOnce.Do(func() {
		orders, err := s.repo.GetAllOrders(ctx)
		if err != nil {
			initErr = err
			return
		}
		for _, order := range orders {
			_ = s.Redis.SetOrder(ctx, &order)
		}
	})
	return initErr
}
