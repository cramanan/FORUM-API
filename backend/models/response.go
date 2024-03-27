package models

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func SendResponse(w http.ResponseWriter, r Response) (err error) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(r.StatusCode)
	err = json.NewEncoder(w).Encode(r)
	return err
}
