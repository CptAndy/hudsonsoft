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
	ErrDuplicateProduct    = errors.New("product_name duplicate not allowed...violates unique contraints")
	ErrDuplicateType       = errors.New("type_name duplicate not allowed...violates unique contraints")
	ErrDuplicateReturnType = errors.New("return_type duplicate not allowed...violates unique contraints")

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
	Products interface {
		Create(context.Context, *Product) error
		GetBySalesNum(context.Context, string) (*Product, error)
		Delete(context.Context, string) error
	}
	ProductTypes interface {
		Create(context.Context, *ProductType) error
		GetByTypeID(context.Context, string) (*ProductType, error)
		Delete(context.Context, string) error
	}
	ReturnTypes interface {
		Create(context.Context, *ReturnType) error
		GetByReturnID(context.Context, string) (*ReturnType, error)
		Delete(context.Context, string) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Employees:    &EmployeeStore{db},
		Customers:    &CustomerStore{db},
		Products:     &ProductStore{db},
		ProductTypes: &ProductTypeStore{db},
		ReturnTypes:  &ReturnTypeStore{db},
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
