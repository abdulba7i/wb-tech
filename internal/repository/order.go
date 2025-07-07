package repository

import (
	"context"
	"database/sql"
	"fmt"
	"wb-tech/internal/model"
)

type OrderRepository interface {
	CreateOrder(ordr model.Order) error
	GetOrderByUID(uid string) (*model.Order, error)
	GetAllOrders(ctx context.Context) ([]model.Order, error)
	GetAllBasicOrders(ctx context.Context) ([]model.Order, error)
	getOrderFromDB(orderUID string) (*model.Order, int64, int64, error)
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &Storage{db: db}
}

func (s *Storage) CreateOrder(ordr model.Order) error {
	const op = "repository.postgres.CreateOrder"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)", ordr.OrderUID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: error check order: %w", op, err)
	}
	if exists {
		return fmt.Errorf("%s: order with UID %s exists", op, ordr.OrderUID)
	}

	idDelivery, err := s.AddDelivery(tx, ordr.Delivery)
	if err != nil {
		return fmt.Errorf("%s: failed to add delivery: %w", op, err)
	}

	idPayment, err := s.AddPayment(tx, ordr.Payment)
	if err != nil {
		return fmt.Errorf("%s: failed to add payment: %w", op, err)
	}

	// query := `INSERT INTO orders
	// 	(order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, oof_shard)
	// 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	// `

	// _, err = tx.Exec(query, ordr.OrderUID, ordr.TrackNumber, ordr.Entry, idDelivery, idPayment, ordr.Locale,
	// 	ordr.InternalSignature, ordr.CustomerID, ordr.DeliveryService, ordr.ShardKey, ordr.SMID, ordr.OOFShard)

	query := `INSERT INTO orders
	(order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, 
	delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err = tx.Exec(query,
		ordr.OrderUID,
		ordr.TrackNumber,
		ordr.Entry,
		idDelivery,
		idPayment,
		ordr.Locale,
		ordr.InternalSignature,
		ordr.CustomerID,
		ordr.DeliveryService,
		ordr.ShardKey,
		ordr.SMID,
		ordr.DateCreated, // <-- добавлено!
		ordr.OOFShard,
	)

	if err != nil {
		return fmt.Errorf("%s: failed to insert order: %w", op, err)
	}

	err = s.AddItems(tx, ordr.OrderUID, ordr.Items)
	if err != nil {
		return fmt.Errorf("%s: failed to add items: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) GetOrderByUID(uid string) (*model.Order, error) {
	const op = "repository.postgres.GetOrderByUID"

	order, deliveryID, paymentID, err := s.getOrderFromDB(uid)
	if err != nil {
		return nil, err
	}

	delivery, err := s.GetDeliveryByID(nil, deliveryID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get delivery: %w", op, err)
	}
	order.Delivery = *delivery

	payment, err := s.GetPaymentByID(nil, paymentID)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get payment: %w", op, err)
	}
	order.Payment = *payment

	items, err := s.GetItemsByOrderUID(nil, uid)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get items: %w", op, err)
	}
	order.Items = items

	return order, nil
}

func (s *Storage) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	const op = "repository.postgres.GetAllOrders"

	query := `
        SELECT 
            order_uid, track_number, entry, 
            delivery_id, payment_id, locale, 
            internal_signature, customer_id, 
            delivery_service, shardkey, sm_id, oof_shard, date_created
        FROM orders
    `

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		var deliveryID, paymentID int64

		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&deliveryID,
			&paymentID,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SMID,
			&order.OOFShard,
			&order.DateCreated,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		delivery, err := s.GetDeliveryByID(nil, deliveryID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get delivery: %w", op, err)
		}
		order.Delivery = *delivery

		payment, err := s.GetPaymentByID(nil, paymentID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get payment: %w", op, err)
		}
		order.Payment = *payment

		items, err := s.GetItemsByOrderUID(nil, order.OrderUID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get items: %w", op, err)
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *Storage) GetAllBasicOrders(ctx context.Context) ([]model.Order, error) {
	const op = "repository.postgres.getAllBasicOrders"

	query := `SELECT order_uid FROM orders`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.OrderUID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *Storage) getOrderFromDB(orderUID string) (*model.Order, int64, int64, error) {
	const op = "repository.postgres.getOrderFromDB"

	query := `
        SELECT 
            order_uid, track_number, entry, 
            delivery_id, payment_id, locale, 
            internal_signature, customer_id, 
            delivery_service, shardkey, sm_id, oof_shard, date_created
        FROM orders 
        WHERE order_uid = $1
    `

	var order model.Order
	var deliveryID, paymentID int64

	err := s.db.QueryRow(query, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&deliveryID,
		&paymentID,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SMID,
		&order.OOFShard,
		&order.DateCreated,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, 0, 0, fmt.Errorf("%s: order not found", op)
		}
		return nil, 0, 0, fmt.Errorf("%s: failed to get order: %w", op, err)
	}

	return &order, deliveryID, paymentID, nil
}
