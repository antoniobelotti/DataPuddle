package main

import (
	"net/http"
)

func SessionKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := GetNewApiKey()
	respondWithJSON(w, http.StatusOK, SessionKeyReponse{Outcome: "ok", Key: key})
}
