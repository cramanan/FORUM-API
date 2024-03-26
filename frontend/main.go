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

	mux.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("static/js"))))
	mux.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))

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
