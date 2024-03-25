package utils

import (
	"net/http"
)

var SERVERMESSAGES = map[int][]byte{
	http.StatusInternalServerError: []byte("Internal Server Error"),
	http.StatusBadRequest:          []byte("Bad Request"),
	http.StatusOK:                  []byte("OK"),
}
