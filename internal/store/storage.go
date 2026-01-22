package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound            = errors.New("record not found")
	ErrConflict            = errors.New("resource already exists...duplicate key values violate constraint")
	ErrDuplicateEmail      = errors.New("email duplicate not allowed...violates unique contraints")
	ErrDuplicateUsername   = errors.New("username duplicate not allowed...violates unique contraints")
	ErrDuplicateEmployeeID = errors.New("employee_id duplicate not allowed...violates unique contraints")
	ErrDuplicateCustomerID = errors.New("customer_id duplicate not allowed...violates unique contraints")

	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Employees interface {
		Create(context.Context, *Employee) error
		GetByID(context.Context, string) (*Employee, error)
		Delete(context.Context, string) error
	}
	Customers interface {
		Create(context.Context, *Customer) error
		GetByID(context.Context, string) (*Customer, error)
		Delete(context.Context, string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Employees: &EmployeeStore{db},
		Customers: &CustomerStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
