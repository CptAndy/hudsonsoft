package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type ReturnType struct {
	ID          int64  `json:"id"`
	Return_Name string `json:"return_name"`
}

type ReturnTypeStore struct {
	db *sql.DB
}

//Create a reason with transaction

func (s *ReturnTypeStore) Create(ctx context.Context, returnType *ReturnType) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.create(ctx, tx, returnType)
	})
}

func (s *ReturnTypeStore) Delete(ctx context.Context, returntypeID string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, returntypeID); err != nil {
			return err
		}
		return nil
	})
}

func (s *ReturnTypeStore) GetByReturnID(ctx context.Context, returntypeID string) (*ReturnType, error){
	query := `SELECT id, return_name FROM return_types WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	returnType := &ReturnType{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		returntypeID,
	).Scan(
		&returnType.ID,
		&returnType.Return_Name,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return returnType, nil
}

func (s *ReturnTypeStore) create(ctx context.Context, tx *sql.Tx, returntype *ReturnType) error {
	query := `INSERT INTO return_types (return_name) VALUES ($1);`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		returntype.Return_Name,
	).Scan(
		&returntype.ID,
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

func (s *ReturnTypeStore) delete(ctx context.Context, tx *sql.Tx, returnTypeID string) error {
	query := `DELETE FROM return_types WHERE type_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, returnTypeID)
	if err != nil {
		return nil
	}
	return nil
}
