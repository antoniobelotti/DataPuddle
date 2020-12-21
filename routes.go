package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func SetUpRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/sessionkey", SessionKeyHandler).
		Methods("GET")

	router.HandleFunc("/pwd", PWDHandler).
		Methods("GET").
		Queries("key", "{key}")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}
