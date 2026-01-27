package store

import (
	"context"
	"database/sql"
)

type ProductType struct {
	ID       int64  `json:"id"`
	Type_ID   string `json:"type_id"`
	Type_Name string `json:"type_name"`
}

type ProductTypeStore struct {
	db *sql.DB
}

// Create product_type with transaction
func (s *ProductTypeStore) Create(ctx context.Context, prodtype *ProductType) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.create(ctx, tx, prodtype)
	})

}

// Create a variation
func (s *ProductTypeStore) create(ctx context.Context, tx *sql.Tx, prodtype *ProductType) error {
	query := `INSERT INTO product_type (type_name) VALUES ($1);`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		prodtype.Type_Name,
	).Scan(
		&prodtype.ID,
	)

	if err != nil {
		switch {
		case err.Error() == `duplicate key value violates unique constraint "products_type_id_key"`:
			return ErrDuplicateTypeID
		default:
			return err
		}
	}
	return nil
}
