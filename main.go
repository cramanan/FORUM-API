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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Frontend setup
		http.ServeFile(w, r, "static/index.html")
	})
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.Handle("/api/", /* API Endpoint*/
		routes.Protected( /* Auth middleware */
			routes.HandleFunc(routes.Root))) /* Route handler*/

	mux.HandleFunc("/api/register",
		routes.HandleFunc(routes.Register)) /*This route is not protected*/

	mux.HandleFunc("/api/login",
		routes.HandleFunc(routes.Login)) /*This route is not protected*/

	mux.Handle("/api/logout",
		routes.Protected(
			routes.HandleFunc(routes.Logout)))

	mux.Handle("/api/getposts",
		routes.Protected(
			routes.HandleFunc(routes.GetPosts)))

	mux.Handle("/api/getusers",
		routes.Protected(
			routes.HandleFunc(routes.GetUsers)))

	mux.Handle("/api/post",
		routes.Protected(
			routes.HandleFunc(routes.Post)))

	// mux.Handle("/api/ws", middleware.Protected(routes.HandleFunc(routes.WS)))
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Server up and running on port %s\n", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error here")
	}
}
