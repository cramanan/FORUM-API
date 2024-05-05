package main

import (
	"backend/database"
	"backend/routes"
	"log"
	"net/http"
)

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
	mux.HandleFunc("/getposts", routes.GetPosts)
	// mux.HandleFunc("/post", routes.Post)
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
