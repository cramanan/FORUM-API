package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"
	"real-time-forum/api/database"
	"real-time-forum/api/models"
	"strconv"
	"time"
)

type API struct {
	http.Server
	Storage  *database.Sqlite3Store
	Sessions *database.SessionStore
}

func NewAPI(addr string) (*API, error) {
	server := new(API)
	server.Server.Addr = addr

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Frontend setup
		http.ServeFile(w, r, "static/index.html")
	})
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandleFunc("/api/", server.Protected(HandleFunc(server.ReadSession)))

	router.HandleFunc("/api/register", HandleFunc(server.Register))
	router.HandleFunc("/api/login", HandleFunc(server.Login))

	router.HandleFunc("/api/users", server.Protected(HandleFunc(server.GetUsers)))
	router.HandleFunc("/api/post", server.Protected(HandleFunc(server.Post)))
	router.HandleFunc("/api/posts", server.Protected(HandleFunc(server.GetPosts)))

	server.Server.Handler = router

	storage, err := database.NewSqlite3Store()
	if err != nil {
		return nil, err
	}

	server.Storage = storage
	server.Sessions = database.NewSessionStore()
	return server, nil
}

func writeJSON(writer http.ResponseWriter, statusCode int, v any) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	return json.NewEncoder(writer).Encode(v)
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func HandleFunc(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func (server *API) Protected(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := server.Sessions.GetSession(r)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (server *API) Register(writer http.ResponseWriter, request *http.Request) error {
	if request.Method != http.MethodPost {
		return writeJSON(writer, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	registerReq := new(models.RegisterRequest)
	err := json.NewDecoder(request.Body).Decode(registerReq)
	if err != nil {
		return err
	}

	if registerReq.Email == "" ||
		registerReq.Name == "" ||
		registerReq.Password == "" ||
		registerReq.Gender == "" ||
		registerReq.Age == "" ||
		registerReq.FirstName == "" ||
		registerReq.LastName == "" {

		return writeJSON(writer, http.StatusBadRequest, "Missing Credentials")
	}

	if _, err = mail.ParseAddress(registerReq.Email); err != nil {
		writeJSON(writer, http.StatusBadRequest, "Invalid Email")
	}

	if _, err = strconv.Atoi(registerReq.Age); err != nil {
		writeJSON(writer, http.StatusBadRequest, "Age field is invalid")
	}

	registerReq.Ctx, registerReq.CancelCtx = context.WithTimeout(request.Context(), 3*time.Second)

	user, err := server.Storage.RegisterUser(registerReq)
	if err != nil {
		return err
	}

	session := server.Sessions.NewSession(writer, request)
	session.User = user

	return writeJSON(writer, http.StatusCreated, user)
}

func (server *API) Login(writer http.ResponseWriter, request *http.Request) error {
	if request.Method != http.MethodPost {
		return writeJSON(writer, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	loginReq := new(models.LoginRequest)
	err := json.NewDecoder(request.Body).Decode(loginReq)
	if err != nil {
		return err
	}

	if loginReq.Email == "" ||
		loginReq.Password == "" {

		return writeJSON(writer, http.StatusBadRequest, "Missing Credentials")
	}

	loginReq.Ctx, loginReq.CancelCtx = context.WithTimeout(request.Context(), 3*time.Second)
	user, err := server.Storage.LogUser(loginReq)
	if err != nil {
		return writeJSON(writer, http.StatusBadRequest, "Invalid Password")
	}

	session := server.Sessions.NewSession(writer, request)
	session.User = user

	return writeJSON(writer, http.StatusOK, user)
}

func (server *API) GetUsers(writer http.ResponseWriter, request *http.Request) error {
	users, err := server.Storage.GetUsers()
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, users)
}

func (server *API) ReadSession(writer http.ResponseWriter, request *http.Request) error {
	session, err := server.Sessions.GetSession(request)
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, session.User)
}

func (server *API) Post(writer http.ResponseWriter, request *http.Request) error {
	if request.Method != http.MethodPost {
		return writeJSON(writer, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	postReq := new(models.PostRequest)
	err := json.NewDecoder(request.Body).Decode(postReq)
	if err != nil {
		return err
	}

	if postReq.Content == "" {
		return writeJSON(writer, http.StatusBadRequest, "Missing Credentials")
	}

	session, err := server.Sessions.GetSession(request)
	if err != nil {
		return err
	}

	postReq.UserID = session.User.ID
	postReq.Username = session.User.Name
	var cancel context.CancelFunc
	postReq.Ctx, cancel = context.WithTimeout(request.Context(), 3*time.Second)
	defer cancel()

	post, err := server.Storage.CreatePost(postReq)
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, post)
}

func (server *API) GetPosts(writer http.ResponseWriter, request *http.Request) error {
	posts, err := server.Storage.GetPosts()
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, posts)
}
