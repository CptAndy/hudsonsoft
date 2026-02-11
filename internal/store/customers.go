package store

import (
	"context"
	"database/sql"
	"strings"
)

type Customer struct {
	ID               int64   `json:"id"`
	Customer_id      string  `json:"customer_id" validate:"required,min=10,max=10"`
	First_name       string  `json:"first_name" validate:"required,max=50"`
	Last_name        string  `json:"last_name" validate:"required,max=50"`
	Email            string  `json:"email" validate:"email"`
	City             string  `json:"city" validate:"required,min=5,max=100"`
	State            string  `json:"state" validate:"required,min=2,max=2"`
	Amount_spent     float64 `json:"amount_spent"`
	Product_owned    int64   `json:"product_owned"`
	Product_returned int64   `json:"product_returned"`
}

// Customer storage DB

type CustomerStore struct {
	db *sql.DB
}

// Find a customer using their ID
func (s *CustomerStore) GetByID(ctx context.Context, custID string) (*Customer, error) {
	query := `SELECT c.cust_id, c.first_name, c.last_name, c.email, c.city, c.state 
			FROM customers c
			WHERE id = (SELECT id 
						FROM customers 
						WHERE cust_id = $1)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	customer := &Customer{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		custID,
	).Scan(
		&customer.Customer_id,
		&customer.First_name,
		&customer.Last_name,
		&customer.Email,
		&customer.City,
		&customer.State,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return customer, nil
}

// Create a customer
func (s *CustomerStore) Create(ctx context.Context, customer *Customer) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.create(ctx, tx, customer)
	})

}

// Delete a customer using their customer_id
func (s *CustomerStore) Delete(ctx context.Context, custID string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, custID); err != nil {
			return err
		}
		return nil
	})
}

// DELETE QUERY
func (s *CustomerStore) delete(ctx context.Context, tx *sql.Tx, custID string) error {
	query := `DELETE FROM customers WHERE cust_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, custID)
	if err != nil {
		return nil
	}
	return nil
}

// CREATE QUERY
func (s *CustomerStore) create(ctx context.Context, tx *sql.Tx, customer *Customer) error {
	query := `
		INSERT INTO customers (first_name,last_name,email,city,state) 
		VALUES ($1,$2,$3,$4,$5)
		returning id
		`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		strings.ToUpper(customer.First_name),
		strings.ToUpper(customer.Last_name),
		customer.Email,
		customer.City,
		customer.State,
	).Scan(
		&customer.ID,
	)
	if err != nil {
		switch {
		case err.Error() == `duplicate key value violates unique constraint "customers_cust_id_key"`:
			return ErrDuplicateCustomerID
		default:
			return err
		}
	}

	return nil
}

// a handler that can search customers by Customer_id or Email or First_name/Last_name
