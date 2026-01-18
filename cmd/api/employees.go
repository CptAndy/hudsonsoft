package main

import (
	"errors"
	"net/http"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
	"github.com/go-chi/chi/v5"
)

var (
	missingID = errors.New("emp_id is missing")
)

type CreateEmployeePayload struct {
	First_name string `json:"first_name" validate:"required,max=100"`
	Last_name  string `json:"last_name" validate:"required,max=100"`
	Password   string `json:"password" validate:"required,min=3,max=72"`
}

type empKey string

const empCtx empKey = "employee"

func (app *application) createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateEmployeePayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	employee := &store.Employee{
		Fname: payload.First_name,
		Lname: payload.Last_name,
	}

	if err := employee.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	// Create a transaction
	err := app.store.Employees.Create(ctx, employee)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	empID := chi.URLParam(r, "employeeID")
	if empID == "" {
		app.badRequestResponse(w, r, missingID)
		return
	}

	ctx := r.Context()

	employee, err := app.store.Employees.GetByID(ctx, empID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, employee); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	empID := chi.URLParam(r, "employeeID")

	ctx := r.Context()

	err := app.store.Employees.Delete(ctx, empID)
	//If theres an error
	if err != nil {
		// go into a switch
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getEmployeeFromCtx(r *http.Request) *store.Employee {
	employee, _ := r.Context().Value(empCtx).(*store.Employee)
	return employee
}
