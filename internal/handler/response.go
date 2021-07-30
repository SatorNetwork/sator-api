package handler

import (
	"encoding/json"
	"net/http"
)

type data map[string]interface{}

// Serve data as JSON as response
func jsonResponse(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Serve data as string as response
func stringResponse(w http.ResponseWriter, status int, str string) error {
	w.Header().Set("Content-Type", "plain/text; charset=UTF-8")
	w.WriteHeader(status)
	_, err := w.Write([]byte(str))
	return err
}
