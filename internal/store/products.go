package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Product struct {
	ID           int64  `json:"id"`
	Product_name string `json:"product_name" validate:"required,min=4,max=50"`
	Sales_num    string `json:"sales_num"`
}

type ProductStore struct {
	db *sql.DB
}

// Find product by sales_num
func (s *ProductStore) GetBySalesNum(ctx context.Context, salesNum string) (*Product, error) {
	query := `
SELECT
    p.sales_num,
    p.product_name
FROM
    products p
WHERE
    id = (
        SELECT
            id
        FROM
            products
        WHERE
            sales_num = $1)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	product := &Product{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		salesNum,
	).Scan(
		&product.Sales_num,
		&product.Product_name,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return product, nil
}

// Create product with transaction
func (s *ProductStore) Create(ctx context.Context, product *Product) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.create(ctx, tx, product)
	})

}

// Create a product
func (s *ProductStore) create(ctx context.Context, tx *sql.Tx, product *Product) error {
	query := `INSERT INTO products (product_name)
    VALUES ($1)
	returning id`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		product.Product_name,
	).Scan(
		&product.ID,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return ErrDuplicateProduct
			default:
				return err
			}
		}

	}
	return nil

}

// Delete with a transaction
func (s *ProductStore) Delete(ctx context.Context, prodID string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, prodID); err != nil {
			return err
		}
		return nil
	})
}

// DELETE QUERY
func (s *ProductStore) delete(ctx context.Context, tx *sql.Tx, prodID string) error {
	query := `DELETE FROM products WHERE sales_num = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, prodID)
	if err != nil {
		return nil
	}
	return nil

}
