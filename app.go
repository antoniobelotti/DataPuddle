package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type App struct {
	router *mux.Router
}

func (app *App) Run(address string) {
	http.ListenAndServe(address, app.router)
}
