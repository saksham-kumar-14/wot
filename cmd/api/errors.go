package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusInternalServerError, "an error occured ;-;")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("not found error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflict(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("conflict error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusConflict, "conflict occured")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retry string) {
	app.logger.Warnf("conflict error", "method", r.Method, "path", r.URL.Path, "error", "rate limit exceeded. Try after : "+retry)
	writeJSONError(w, http.StatusConflict, "conflict occured")
}
