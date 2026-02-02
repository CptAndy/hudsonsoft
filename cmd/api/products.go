package main

import (
	"net/http"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateProductPayload struct {
	Product_name string `json:"product_name" validate:"required,min=2,max=50"`
}
type prodKey string

const prodCtx prodKey = "product"

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateProductPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	product := &store.Product{
		Product_name: payload.Product_name,
	}
	ctx := r.Context()

	err := app.store.Products.Create(ctx, product)
	if err != nil {
		switch err {
		case store.ErrDuplicateProduct:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}

	}
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	prodID := chi.URLParam(r, "productID")

	if prodID == "" {
		app.badRequestResponse(w, r, missingCustID)
		return
	}

	ctx := r.Context()

	product, err := app.store.Products.GetBySalesNum(ctx, prodID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, product); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	prodID := chi.URLParam(r, "productID")

	ctx := r.Context()

	err := app.store.Products.Delete(ctx, prodID)
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
