package main

import (
	"net/http"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
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
		app.internalServerError(w, r, err)
		return
	}

}
