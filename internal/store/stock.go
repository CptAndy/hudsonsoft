package store

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
)

type Stock struct {
	ID           int64   `json:"id"`
	Product_name string  `json:"product_name"`
	Product_ID   string  `json:"product_id"`
	Type_ID      string  `json:"type_id"`
	Stock        int64   `json:"stock"`
	Price        float64 `json:"price"`
	InStock      bool    `json:"inStock"`
	OnOrder      bool    `json:"onOrder"`
}

type StockStore struct {
	db *sql.DB
}

// Create wrapped in a Transaction
func (s *StockStore) Create(ctx context.Context, stock *Stock) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.create(ctx, tx, stock)
	})
}

func (s *StockStore) Delete(ctx context.Context, product_id string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.delete(ctx, tx, product_id)
	})
}
func (s *StockStore) GetByID(ctx context.Context, product_id string) (*Stock, error) {
	query := `SELECT s.id, s.product_name, s.product_id, s.stock, s.price, s.instock, s.onorder FROM stock s WHERE s.product_id = $1`

ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
defer cancel()

stock := &Stock{}

err := s.db.QueryRowContext(
	ctx,
	query,
	product_id,
).Scan(
	&stock.ID,
	&stock.Product_name,
	&stock.Product_ID,
	&stock.Stock,
	&stock.Price,
	&stock.InStock,
	&stock.OnOrder,
)

if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return stock, nil

}

func (s *StockStore) create(ctx context.Context, tx *sql.Tx, stock *Stock) error {
	log.Printf("INSIDE create------: %s and %s", stock.Product_ID, stock.Type_ID)

	query := `INSERT INTO stock (product_name, product_id, price)
SELECT
    p.product_name || ' ' || pt.type_name,
    p.sales_num || '' || pt.type_id,
	$3
FROM
    products p
    JOIN product_types pt ON TRUE
WHERE
    p.sales_num = $1
    AND pt.type_id = $2
RETURNING ID id;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		stock.Product_ID,
		stock.Type_ID,
		stock.Price,
	).Scan(
		&stock.ID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return ErrDuplicateStock
			default:
				return err
			}
		}
	}
	return nil

}

func (s *StockStore) delete(ctx context.Context, tx *sql.Tx, product_id string) error {
	query := `DELETE FROM stock WHERE product_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, product_id)
	if err != nil {
		return nil
	}
	return nil

}
