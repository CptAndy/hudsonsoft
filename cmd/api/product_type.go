package main

import (
	"net/http"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateProductTypePayload struct {
	Type_name string `json:"type_name" validate:"required,min=2,max=15"`
}

type typeKey string

const typectx typeKey = "product_type"

func (app *application) createProductTypeHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateProductTypePayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product_type := &store.ProductType{
		Type_Name: payload.Type_name,
	}
	ctx := r.Context()

	err := app.store.ProductTypes.Create(ctx, product_type)
	if err != nil {
		app.conflictResponse(w, r, err)
		return
	}

}

func (app *application) getProductTypeHandler(w http.ResponseWriter, r *http.Request) {
	prodTypeID := chi.URLParam(r, "prodTypeID")

	ctx := r.Context()

	product_type, err := app.store.ProductTypes.GetByTypeID(ctx, prodTypeID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, product_type); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) deleteProductTypeHandler(w http.ResponseWriter, r *http.Request) {
	prodTypeID := chi.URLParam(r, "prodTypeID")
	ctx := r.Context()

	err := app.store.ProductTypes.Delete(ctx, prodTypeID)
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
