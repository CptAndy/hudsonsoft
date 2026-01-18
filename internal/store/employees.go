package store

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Employee struct {
	ID         int64    `json:"id"`
	Emp_id     string   `json:"emp_id"`
	First_name string   `json:"first_name"`
	Last_name  string   `json:"last_name"`
	Password   password `json:"-"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	//Encryption
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil

}

type EmployeeStore struct {
	db *sql.DB
}

func (s *EmployeeStore) Create(ctx context.Context, employee *Employee) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		return s.create(ctx, tx, employee)
	})

}

func (s *EmployeeStore) create(ctx context.Context, tx *sql.Tx, employee *Employee) error {
	query := `
	INSERT INTO employees (first_name, last_name, employee_pass)
VALUES ($1,$2,$3) returning id
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(
		ctx,
		query,
		employee.First_name,
		employee.Last_name,
		employee.Password.hash,
	).Scan(

		&employee.ID,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique contraints "employees_emp_id_key"`:
			return ErrDuplicateEmployeeID
		default:
			return err
		}
	}
	return nil
}

func (s *EmployeeStore) GetByID(ctx context.Context, empID string) (*Employee, error) {
	query := ` 
		SELECT e.id, e.emp_id, e.first_name, e.last_name, e.employee_pass
		FROM employees e
		WHERE e.emp_id = $1 
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	employee := &Employee{}
	var passwordHash []byte
	err := s.db.QueryRowContext(
		ctx,
		query,
		empID,
	).Scan(
		&employee.ID,
		&employee.Emp_id,
		&employee.First_name,
		&employee.Last_name,
		&passwordHash,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err

		}
	}
	return employee, nil
}

func (s *EmployeeStore) Delete(ctx context.Context, empID string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.delete(ctx, tx, empID); err != nil {
			return err
		}

		return nil

	})
}

func (s *EmployeeStore) delete(ctx context.Context, tx *sql.Tx, empID string) error {
	query := `DELETE FROM employees WHERE emp_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, empID)
	if err != nil {
		return nil
	}
	return nil
}

func (p *password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}
