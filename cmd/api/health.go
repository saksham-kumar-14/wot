package main

import "net/http"

func (app *application) healthChecker(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))

}
