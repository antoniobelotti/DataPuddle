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

	router.HandleFunc("/cd", CDHandler).
		Methods("GET").
		Queries("key", "{key}", "path", "{path}")

	router.HandleFunc("/mkdir", MKDIRHandler).
		Methods("GET").
		Queries("key", "{key}", "path", "{path}")

	router.HandleFunc("/store", STOREHandler).
		Methods("POST").
		Queries("key", "{key}", "filename", "{filename}")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}
