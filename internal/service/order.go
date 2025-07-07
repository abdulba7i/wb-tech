package service

import (
	"context"
	"wb-tech/internal/cache"
	"wb-tech/internal/model"
	"wb-tech/internal/repository"
)

type OrderService struct {
	repo  *repository.Repository
	cache *cache.CacheService
}

func NewOrderService(repo *repository.Repository, cache *cache.CacheService) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderService) GetOrderByUID(ctx context.Context, uid string) (*model.Order, error) {
	return s.cache.Get(ctx, uid)
}

func (s *OrderService) CreateOrder(ctx context.Context, order *model.Order) error {
	if err := s.repo.Order.CreateOrder(*order); err != nil {
		return err
	}
	return s.cache.Set(ctx, order)
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	return s.repo.Order.GetAllBasicOrders(ctx)
}

func (s *OrderService) PreloadCache(ctx context.Context) error {
	return s.cache.Preload(ctx)
}
