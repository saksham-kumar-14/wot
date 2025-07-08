package main

import "net/http"

func (app *application) healthChecker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{ "status" : "ok" }`))
}
