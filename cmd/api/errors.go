package main

import (
	"net/http"
)

// Handles errors sent via HTTP/API

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJsonError(w, http.StatusInternalServerError, "Server encountered an issue!")

}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {

	app.logger.Warnw("forbidden response", "method", r.Method, "path", r.URL.Path, "error")

	WriteJsonError(w, http.StatusForbidden, "forbidden")

}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJsonError(w, http.StatusBadRequest, err.Error())

}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJsonError(w, http.StatusConflict, err.Error())

}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnf("not found", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	WriteJsonError(w, http.StatusNotFound, "not found")

}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	WriteJsonError(w, http.StatusUnauthorized, "unauthorized")

}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized basic", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	WriteJsonError(w, http.StatusUnauthorized, "unauthorized")

}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnf("unauthorized", "method", r.Method, "path", r.URL.Path)
	WriteJsonError(w, http.StatusUnauthorized, "unauthorized")

}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)
	WriteJsonError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
