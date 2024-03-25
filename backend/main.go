package main

import (
	"backend/routes"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", routes.BasicUpgrade)
	mux.HandleFunc("/register", routes.RegisterClient)
	server := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
