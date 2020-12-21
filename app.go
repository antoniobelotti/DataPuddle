package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

type App struct {
	router *mux.Router
}

func (app *App) Run(address string)  {
	http.ListenAndServe(address, app.router)
}
