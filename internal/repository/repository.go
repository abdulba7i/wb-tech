package repository

import (
	"database/sql"
	"wb-tech/internal/model"
)

type Order interface {
	CreateOrder(ordr model.Order) error
	// GetOrderById(id string) (*model.Order, error)
	// getOrderMainInfo(orderUID string) (*model.Order, error)
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
