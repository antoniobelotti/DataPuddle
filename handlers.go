package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
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

func STOREHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["key"]
	currentPath := sessions.Get(key)
	if currentPath == "" {
		respondWithJSON(w, http.StatusBadRequest, OutcomeResponse{Outcome: "error"})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, OutcomeResponse{Outcome: "error"})
		return
	}

	absFilePath := filepath.Join(actualPath(currentPath), vars["filename"])

	if _, err := os.Stat(absFilePath); err == nil || os.IsExist(err) {
		// file already exists
		respondWithJSON(w, http.StatusBadRequest, OutcomeResponse{Outcome: "error"})
		return
	}

	if err = ioutil.WriteFile(absFilePath, body, 0777); err != nil {
		respondWithJSON(w, http.StatusBadRequest, OutcomeResponse{Outcome: "error"})
		return
	}
	respondWithJSON(w, http.StatusOK, OutcomeResponse{Outcome: "ok"})
}

func RETRIEVEHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["key"]
	currentPath := sessions.Get(key)
	if currentPath == "" {
		respondWithJSON(w, http.StatusBadRequest, OutcomeResponse{Outcome: "error"})
		return
	}

	absFilePath := actualPath(filepath.Join(currentPath, vars["filename"]))

	fileContent, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, RetrieveResponse{Outcome: "error"})
	} else {
		respondWithJSON(w, http.StatusOK, RetrieveResponse{Outcome: "ok", File: string(fileContent)})
	}
}
