package repository

import (
	"database/sql"
	"fmt"
	"wb-tech/internal/model"
)

type ItemsRepository interface {
	AddItems(tx *sql.Tx, order_uid string, items []model.Item) error
	GetItemsByOrderUID(tx *sql.Tx, orderUID string) ([]model.Item, error)
}

func NewItemsRepository(db *sql.DB) ItemsRepository {
	return &Storage{db: db}
}

func (s *Storage) AddItems(tx *sql.Tx, order_uid string, items []model.Item) error {
	const op = "repository.postgres.AddItems"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: error to begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	query := `INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, 
		size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	for _, item := range items {
		_, err := tx.Exec(query, order_uid, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("%s: error to insert item: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: error to commit transaction: %w", op, err)
	}

	return nil
}

func (s *Storage) GetItemsByOrderUID(tx *sql.Tx, orderUID string) ([]model.Item, error) {
	const op = "repository.postgres.GetItemsByOrderUID"

	query := `SELECT * FROM items WHERE order_uid = $1`
	rows, err := tx.Query(query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NMID,
			&item.Brand,
			&item.Status,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: failed to get model delivery: %w", op, err)
		}
		items = append(items, item)
	}

	return items, nil
}
