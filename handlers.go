package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
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

func CDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["key"]
	currentPath := sessions.Get(key)
	if currentPath == "" {
		respondWithJSON(w, http.StatusBadRequest, OutcomeResponse{Outcome: "error"})
		return
	}

	var resultingPath string
	receivedPath := vars["path"]
	switch receivedPath {
	case "..":
		resultingPath = filepath.Dir(currentPath)
	case "/":
		resultingPath = "/"
	default:
		resultingPath = filepath.Join(currentPath, receivedPath)
	}

	if _, err := os.Stat(actualPath(resultingPath)); os.IsNotExist(err) {
		// dir does not exist
		respondWithJSON(w, http.StatusBadRequest, OutcomeResponse{Outcome: "error"})
		return
	}

	sessions.Add(key, resultingPath)
	respondWithJSON(w, http.StatusOK, OutcomeResponse{Outcome: "ok"})
}

func MKDIRHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["key"]
	currentPath := sessions.Get(key)
	if currentPath == "" {
		respondWithJSON(w, http.StatusBadRequest, OutcomeResponse{Outcome: "error"})
		return
	}

	newDirPath := vars["path"]
	actualPath := actualPath(filepath.Join(currentPath, newDirPath))
	os.MkdirAll(actualPath, 0777)
	respondWithJSON(w, http.StatusOK, OutcomeResponse{Outcome: "ok"})
}