package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/execute", app.createExecTask).Methods(http.MethodPost)

	return router
}
