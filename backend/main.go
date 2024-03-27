package main

import (
	"backend/database"
	"backend/models"
	"backend/routes"
	"backend/utils"
	"log"
	"net/http"
)

func init() {
	utils.Store = models.NewStore()
}

func main() {
	err := database.InitDB()
	if err != nil {
		log.Println(err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", routes.Root)
	mux.HandleFunc("/register", routes.RegisterClient)
	mux.HandleFunc("/login", routes.LogClientIn)
	server := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	log.Println("Backend Server up...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
