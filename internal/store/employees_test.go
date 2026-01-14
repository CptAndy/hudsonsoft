package store

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

//Global mock DB vars

var (
	db   *sql.DB
	mock sqlmock.Sqlmock
	err  error
)

func TestMain(m *testing.M) {
	log.Println("---Begin Test---")

	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("error opening mock db %s", err)
	}

	defer db.Close()

	code := m.Run()

	log.Println("---End Test---")
	os.Exit(code)
}

func TestEmployeeCreate(t *testing.T) {
	store := &EmployeeStore{db: db}

	employee := &Employee{
		Fname: "Test",
		Lname: "Test",
	}

	err := employee.Password.Set("test")
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO employees`).WithArgs(
		employee.Fname,
		employee.Lname,
		sqlmock.AnyArg(),
	).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(int64(1)),
		)

	mock.ExpectCommit()

	err = store.Create(context.Background(), employee)

	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}

	if employee.ID != 1 {
		t.Fatalf("expected employee ID 1, got %d", employee.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestEmployeeGetById(t *testing.T) {
	store := &EmployeeStore{db: db}

	employee_id := "T123"
	employee := &Employee{
		ID:     1,
		Emp_id: employee_id,
		Fname:  "Test",
		Lname:  "Test",
	}

	err := employee.Password.Set("test")
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectQuery(`SELECT e.id, e.emp_id, e.first_name, e.last_name, e.employee_pass
		FROM employees e`).
		WithArgs(employee_id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "emp_id", "first_name", "last_name", "employee_pass"}).
				AddRow(employee.ID, employee.Emp_id, employee.Fname, employee.Lname, []byte("hash")),
		)

	employee, err = store.GetByID(context.Background(), employee_id)
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	if employee.Emp_id != employee_id {
		t.Fatalf("expected employee_id=T123, got %s", employee.Emp_id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
	// log.Printf("expected employee_id=T123, got employee_id=%s", employee.Emp_id)
}

// func TestEmployeeDelete(t *testing.T) {
// 	store := &EmployeeStore{db: db}

// 	employee := &Employee{
// 		ID: 1,
// 		Emp_id: "T111",
// 		Fname: "Test",
// 		Lname: "Test",
// 	}

// 	mock.ExpectBegin()

// 	mock.ExpectQuery(`DELETE FROM employees`)

// }
