package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

var sessions Sessions

func init() {
	sessions = Sessions{}
}

func SessionKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := GetNewApiKey()
	sessions.Add(key, "/")
	respondWithJSON(w, http.StatusOK, SessionKeyReponse{Outcome: "ok", Key: key})
}

func PWDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["key"]
	path := sessions.Get(key)
	if path == "" {
		respondWithJSON(w, http.StatusBadRequest, PWDResponse{Outcome: "error", Path: ""})
	} else {
		respondWithJSON(w, http.StatusOK, PWDResponse{Outcome: "ok", Path: path})
	}
}
