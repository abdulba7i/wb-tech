package repository

import (
	"database/sql"
	"fmt"
	"wb-tech/internal/model"
)

type OrderRepository interface {
	CreateOrder(ordr model.Order) error
	// GetOrderById(id string) (*model.Order, error)
	// getOrderMainInfo(orderUID string) (*model.Order, error)
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

	query := `INSERT INTO orders
		(order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = tx.Exec(query, ordr.OrderUID, ordr.TrackNumber, ordr.Entry, idDelivery, idPayment, ordr.Locale,
		ordr.InternalSignature, ordr.CustomerID, ordr.DeliveryService, ordr.ShardKey, ordr.SMID, ordr.OOFShard)

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

// func (s *Storage) GetOrderByUID(uid string) (*model.Order, error) {
// 	// 1. Получаем основную информацию о заказе
// 	order, err := s.getOrderMainInfo(uid) // SELECT * FROM orders WHERE order_uid = $1
// 	if err != nil {
// 		return nil, fmt.Errorf("order not found: %w", err)
// 	}

// 	// 2. Получаем delivery, payment, items
// 	delivery, err := s.GetDeliveryByID(nil, order.DeliveryID) // nil, т.к. транзакция не нужна
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get delivery: %w", err)
// 	}

// 	payment, err := s.GetPaymentByID(nil, order.PaymentID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get payment: %w", err)
// 	}

// 	items, err := s.GetItemsByOrderUID(nil, order.OrderUID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get items: %w", err)
// 	}

// 	order.Delivery, order.Payment, order.Items = *delivery, *payment, items

// 	return order, nil
// }

// func (s *Storage) getOrderMainInfo(orderUID string) (*model.Order, error) {
// 	const op = "repository.postgres.getOrderFromDB"

// 	query := `
//         SELECT
//             order_uid, track_number, entry,
//             delivery_id, payment_id, locale,
//             internal_signature, customer_id,
//             delivery_service, shardkey, sm_id, oof_shard, date_created
//         FROM orders
//         WHERE order_uid = $1
//     `

// 	var order model.Order
// 	err := s.db.QueryRow(query, orderUID).Scan(
// 		&order.OrderUID,
// 		&order.TrackNumber,
// 		&order.Entry,
// 		&order.Delivery, // Предполагаем, что в model.Order есть поля DeliveryID и PaymentID
// 		&order.Payment,
// 		&order.Locale,
// 		&order.InternalSignature,
// 		&order.CustomerID,
// 		&order.DeliveryService,
// 		&order.ShardKey,
// 		&order.SMID,
// 		&order.OOFShard,
// 		&order.DateCreated,
// 	)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("%s: order not found", op)
// 		}
// 		return nil, fmt.Errorf("%s: failed to get order: %w", op, err)
// 	}

// 	return &order, nil
// }
