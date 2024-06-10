package server

import (
	"context"
	"fmt"
	"net/http"
	"real-time-backend/backend/database"

	//"real-time-backend/backend/server/contexte"
	"time"

	"github.com/gorilla/mux"
)

func LoadServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/api/users", userHandler).Methods("GET")
	router.HandleFunc("/api/posts", postHandler).Methods("GET")
	router.HandleFunc("/api/categories", categorieHandler).Methods("GET")
	router.HandleFunc("/api/comments", commentsHandler).Methods("GET")
	router.HandleFunc("/api/likescomments", LikesCommentsHandler).Methods("GET")
	router.HandleFunc("/api/postscategories", PostscategoriesHandler).Methods("GET")
	router.HandleFunc("/api/postslikes", PostsLikesHandler).Methods("GET")
	router.HandleFunc("/api/session", SessionHandler).Methods("GET")
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if r.Method == http.MethodPost {
			nickname := r.FormValue("nickname")
			age := r.FormValue("age")
			gender := r.FormValue("gender")
			firstname := r.FormValue("firstname")
			lastname := r.FormValue("lastname")
			email := r.FormValue("email")
			password := r.FormValue("password")

			err := database.RegisterUser(ctx, nickname, age, gender, firstname, lastname, email, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "User %s successfully registered", nickname)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}).Methods("POST")

	serverConfig := ServerParameters(router, 10)
	fmt.Println("Server started at 127.0.0.1:8080")
	if err := serverConfig.ListenAndServe(); err != nil {
		fmt.Println("Error to serve:", err)
	}
}
