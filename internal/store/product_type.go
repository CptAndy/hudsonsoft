package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type ProductType struct {
	ID        int64  `json:"id"`
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

// GetByID
func (s *ProductTypeStore) GetByTypeID(ctx context.Context, prodTypeID string) (*ProductType, error) {
	query := `
SELECT
	pt.id,
    pt.type_id,
    pt.type_name
FROM
    product_types pt
WHERE
    id = (
        SELECT
            id
        FROM
            product_types
        WHERE
            type_id = $1)
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	productType := &ProductType{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		prodTypeID,
	).Scan(
		&productType.ID,
		&productType.Type_ID,
		&productType.Type_Name,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return productType, nil
}

// Create a variation
func (s *ProductTypeStore) create(ctx context.Context, tx *sql.Tx, prodtype *ProductType) error {
	query := `INSERT INTO product_types (type_name) VALUES ($1);`

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
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return ErrDuplicateType
			default:
				return err
			}
		}

	}
	return nil
}

// Delete a with a transaction
func (s *ProductTypeStore) Delete(ctx context.Context, prodTypeID string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, prodTypeID); err != nil {
			return err
		}
		return nil
	})

}

func (s *ProductTypeStore) delete(ctx context.Context, tx *sql.Tx, prodTypeID string) error {
	query := `DELETE from product_types WHERE type_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, prodTypeID)
	if err != nil {
		return nil
	}
	return nil
}
