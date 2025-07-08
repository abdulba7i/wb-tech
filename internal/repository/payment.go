package repository

import (
	"database/sql"
	"fmt"
	"wb-tech/internal/model"
)

type PaymentRepository interface {
	AddPayment(tx *sql.Tx, payment model.Payment) (int64, error)
	GetPaymentByID(tx *sql.Tx, id int64) (*model.Payment, error)
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &Storage{db: db}
}

func (s *Storage) AddPayment(tx *sql.Tx, payment model.Payment) (int64, error) {
	const op = "repository.postgres.AddPayment"
	var id int64

	query := `INSERT INTO payment(
				transaction, request_id, currency, provider, amount, 
				payment_dt, bank, delivery_cost, goods_total, custom_fee) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	err := tx.QueryRow(query,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDT,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: error added data for Payment table: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetPaymentByID(tx *sql.Tx, id int64) (*model.Payment, error) {
	const op = "repository.postgres.GetPaymentByID"

	query := `SELECT 
        transaction, request_id, currency, provider, amount,
        payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payment WHERE id = $1`
	var payment model.Payment

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(query, id)
	} else {
		row = s.db.QueryRow(query, id)
	}

	err := row.Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDT,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &payment, nil
}
