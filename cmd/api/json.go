package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate validator.Validate

func init(){
	Validate = *validator.New(validator.WithRequiredStructEnabled())
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	//Headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	//Limit body size to prevent ddos attacks
	//1Mb
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	//decoder
	decoder := json.NewDecoder(r.Body)
	//disallow any unknown field
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJsonError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJson(w, status, &envelope{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return writeJson(w, status, &envelope{Data: data})
}
