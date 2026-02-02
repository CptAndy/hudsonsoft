package main

import (
	"errors"
	"net/http"

	"github.com/CptAndy/hudsonsoftbackend/internal/store"
	"github.com/go-chi/chi/v5"
)

var (
	missingReturnType_id = errors.New("return_type_id is missing")
)

type CreateaReturnTypePayload struct {
	Return_type_name string `json:"return_type_name" validate:"required,max=50"`
}

type returnKey string

const returnCtx returnKey = "return"

func (app *application) createReturnReasonHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateaReturnTypePayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	reasoning := &store.ReturnType{
		Return_Name: payload.Return_type_name,
	}

	ctx := r.Context()

	err := app.store.ReturnTypes.Create(ctx, reasoning)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deleteReturnReasonHandler(w http.ResponseWriter, r *http.Request) {
	returnTypeID := chi.URLParam(r, "returnTypeID")

	ctx := r.Context()

	err := app.store.ReturnTypes.Delete(ctx, returnTypeID)

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
