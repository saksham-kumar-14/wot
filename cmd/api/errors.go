package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s\n", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "an error occured ;-;")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error %s\n", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error %s\n", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflict(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: %s path: %s error %s\n", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusConflict, "conflict occured")
}
