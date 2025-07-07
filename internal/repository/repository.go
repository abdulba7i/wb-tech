package repository

import (
	"context"
	"database/sql"
	"wb-tech/internal/model"
)

type Order interface {
	CreateOrder(ordr model.Order) error
	GetOrderByUID(uid string) (*model.Order, error)
	GetAllOrders(ctx context.Context) ([]model.Order, error)
	GetAllBasicOrders(ctx context.Context) ([]model.Order, error)
	getOrderFromDB(orderUID string) (*model.Order, int64, int64, error)
}

type Payment interface {
	AddPayment(tx *sql.Tx, ayment model.Payment) (int64, error)
	GetPaymentByID(tx *sql.Tx, id int64) (*model.Payment, error)
}

type Items interface {
	AddItems(tx *sql.Tx, order_uid string, items []model.Item) error
	GetItemsByOrderUID(tx *sql.Tx, orderUID string) ([]model.Item, error)
}

type Delivery interface {
	AddDelivery(tx *sql.Tx, delivery model.Delivery) (int64, error)
	GetDeliveryByID(tx *sql.Tx, id int64) (*model.Delivery, error)
}

type Repository struct {
	Order
	Payment
	Items
	Delivery
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Order:    NewOrderRepository(db),
		Payment:  NewPaymentRepository(db),
		Items:    NewItemsRepository(db),
		Delivery: NewDeliveryRepository(db),
	}
}
