package service

import (
	"context"
	"wb-tech/internal/model"
)

type Order interface {
	GetOrderByUID(ctx context.Context, uid string) (*model.Order, error)
	CreateOrder(ctx context.Context, order *model.Order) error
	GetAllOrders(ctx context.Context) ([]model.Order, error)
	PreloadCache(ctx context.Context) error
}

type Service struct {
	Order
}
