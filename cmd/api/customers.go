package main

import (
	"errors"
	"net/http"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
	"github.com/go-chi/chi/v5"
)

var (
	missingCustID = errors.New("emp_id is missing")
)

type CreateCustomerPayload struct {
	First_name string `json:"first_name" validate:"required,max=50"`
	Last_name  string `json:"last_name" validate:"required,max=50"`
	Email      string `json:"email" validate:"email"`
	City       string `json:"city" validate:"required,min=2,max=50"`
	State      string `json:"state" validate:"required,min=2,max=2"`
}

type custKey string

const custCtx custKey = "customer"

func (app *application) createCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCustomerPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	customer := &store.Customer{
		First_name: payload.First_name,
		Last_name:  payload.Last_name,
		Email:      payload.Email,
		City:       payload.City,
		State:      payload.State,
	}
	ctx := r.Context()

	err := app.store.Customers.Create(ctx, customer)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getCustomerHandler(w http.ResponseWriter, r *http.Request) {
	custID := chi.URLParam(r, "customerID")

	if custID == "" {
		app.badRequestResponse(w, r, missingCustID)
		return
	}

	ctx := r.Context()

	customer, err := app.store.Customers.GetByID(ctx, custID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, customer); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) deleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	custID := chi.URLParam(r, "customerID")

	ctx := r.Context()

	err := app.store.Customers.Delete(ctx, custID)

	if err != nil {
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
