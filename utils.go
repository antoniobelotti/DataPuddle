package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"
)

func GetNewApiKey() string {
	key := []byte(time.Now().String())
	hash := md5.Sum(key)
	return hex.EncodeToString(hash[:])
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
