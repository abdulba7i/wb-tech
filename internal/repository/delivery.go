package repository

import (
	"database/sql"
	"fmt"
	"wb-tech/internal/model"
)

type DeliveryRepository interface {
	AddDelivery(tx *sql.Tx, delivery model.Delivery) (int64, error)
	GetDeliveryByID(tx *sql.Tx, id int64) (*model.Delivery, error)
}

func NewDeliveryRepository(db *sql.DB) DeliveryRepository {
	return &Storage{db: db}
}

func (s *Storage) AddDelivery(tx *sql.Tx, delivery model.Delivery) (int64, error) {
	const op = "repository.postgres.AddDelivery"
	var id int64

	query := "INSERT INTO delivery (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	err := tx.QueryRow(query, delivery.Name, delivery.Phone, delivery.Zip, delivery.City,
		delivery.Address, delivery.Region, delivery.Email).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%s: no rows affected", op)
		}
		return 0, fmt.Errorf("%s: error added data for Delivery table: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetDeliveryByID(tx *sql.Tx, id int64) (*model.Delivery, error) {
	const op = "repository.postgres.GetDeliveryByID"

	query := `SELECT name, phone, zip, city, address, region, email FROM delivery WHERE id = $1`
	var delivery model.Delivery

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRow(query, id)
	} else {
		row = s.db.QueryRow(query, id)
	}

	err := row.Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &delivery, nil
}
