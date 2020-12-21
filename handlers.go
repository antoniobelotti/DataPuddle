package main

import (
	"net/http"
	"path/filepath"
)

var sessions Sessions

func init() {
	sessions = Sessions{}
}

func SessionKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := GetNewApiKey()
	sessions.Add(key, filepath.Join("/"))
	respondWithJSON(w, http.StatusOK, SessionKeyReponse{Outcome: "ok", Key: key})
}
