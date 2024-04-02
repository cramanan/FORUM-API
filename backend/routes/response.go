package routes

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int                    `json:"status"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
}

func SendResponse(w http.ResponseWriter, r Response) (err error) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(r.StatusCode)
	return json.NewEncoder(w).Encode(r)
}
