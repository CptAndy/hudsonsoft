package main

import (
	"log"
	"net/http"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
)

type CreateStockPayload struct {
	Product_ID string  `json:"product_id" validate:"required,min=3,max=50"`
	Type_ID    string  `json:"type_id" validate:"required,min=3,max=50"`
	Price      float64 `json:"price"`
}

func (app *application) createStockHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateStockPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	stock := &store.Stock{
		Product_ID: payload.Product_ID,
		Type_ID:    payload.Type_ID,
		Price:      payload.Price,
	}

	err := app.store.Stock.Create(ctx, stock)
	if err != nil {
		//log.Printf("HANDLER ERR TYPE: %T", err)
		//log.Printf("HANDLER ERR VALUE: %+v", err)
		switch err {
		case store.ErrDuplicateStock:
			log.Println("INSIDE DUPLICATE ERROR")
			app.conflictResponse(w, r, err)
			return
		case store.ErrNotFound:
			log.Println("INSIDE NOT FOUND ERROR")
			app.notFoundResponse(w, r, err)
		default:
			log.Println("INSIDE INTERNAL SERVER ERROR RESPONSE")
			app.internalServerError(w, r, err)
			return
		}
	}

}
