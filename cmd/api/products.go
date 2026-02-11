package main

import (
	"log"
	"net/http"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateProductPayload struct {
	Product_name string `json:"product_name" validate:"required,min=2,max=50"`
}


func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateProductPayload

	if err := readJson(w, r, &payload); err != nil {

		log.Println("INSIDE READ JSON")

		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		log.Println("INSIDE VALIDATE JSON")
		app.badRequestResponse(w, r, err)
		return
	}

	product := &store.Product{
		Product_name: payload.Product_name,
	}
	ctx := r.Context()

	err := app.store.Products.Create(ctx, product)
	if err != nil {
		//log.Printf("HANDLER ERR TYPE: %T", err)
		//log.Printf("HANDLER ERR VALUE: %+v", err)
		switch err {
		case store.ErrDuplicateProduct:
			//log.Println("INSIDE CONFLICT RESPONSE")
			app.conflictResponse(w, r, err)
			return
		default:
			//log.Println("INSIDE INTERNAL SERVER ERROR RESPONSE")
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
