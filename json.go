package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, body any) error {
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(body)
}
