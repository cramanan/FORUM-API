package main

import (
	"log"
	"net/http"
	"real-time-forum/api/database"
	"real-time-forum/api/routes"
)

func main() {

	err := database.InitDB()
	if err != nil {
		log.Println(err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/api/", routes.Root)
	mux.HandleFunc("/api/register", routes.RegisterClient)
	mux.HandleFunc("/api/login", routes.LogClientIn)
	mux.HandleFunc("/api/getposts", routes.GetPosts)
	mux.HandleFunc("/api/post", routes.Post)
	// mux.HandleFunc("/logout", routes.Logout)
	mux.HandleFunc("/api/ws", routes.WS)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server up")
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error")
	}
}
