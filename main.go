package main

import (
	"log"
	"net/http"
	"real-time-forum/api/database"
	"real-time-forum/api/routes"
	"real-time-forum/api/routes/middleware"
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

	mux.Handle("/api/", middleware.Protected(routes.Root))
	mux.HandleFunc("/api/register", routes.Register)
	mux.HandleFunc("/api/login", routes.Login)
	mux.Handle("/api/logout", middleware.Protected(routes.Logout))
	mux.Handle("/api/getposts", middleware.Protected(routes.GetPosts))
	mux.Handle("/api/getusers", middleware.Protected(routes.GetUsers))
	mux.Handle("/api/post", middleware.Protected(routes.Post))
	mux.Handle("/api/ws", middleware.Protected(routes.WS))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server up and running...")
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error here")
	}
}
