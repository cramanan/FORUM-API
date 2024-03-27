package main

import (
	"backend/database"
	"backend/models"
	"backend/routes"
	"backend/utils"
	"log"
	"net/http"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Println(err)
		return
	}

	utils.Store = models.NewStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/", routes.BasicUpgrade)
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
