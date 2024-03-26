package utils

import (
	"encoding/json"
	"net/http"
)

var SERVERMESSAGES = map[int][]byte{
	http.StatusInternalServerError: []byte("Internal Server Error"),
	http.StatusBadRequest:          []byte("Bad Request"),
	http.StatusOK:                  []byte("OK"),
	http.StatusUnauthorized:        []byte("Unauthorized"),
}

type Response struct {
	StatusCode int `json:"status"`
	Header     http.Header
	Message    string `json:"message"`
}

func SendResponse(w http.ResponseWriter, r Response) (err error) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(r.StatusCode)
	err = json.NewEncoder(w).Encode(r)
	return err
}
