package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Frontend Server up...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
