package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"
	"real-time-forum/api/database"
	"real-time-forum/api/models"
	"strconv"
)

type API struct {
	http.Server
	Storage  *database.Sqlite3Store
	Sessions *database.SessionStore
	// Upgrader websocket.Upgrader
}

func NewAPI(addr string) (*API, error) {
	server := new(API)
	server.Server.Addr = addr

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Frontend setup
		http.ServeFile(w, r, "static/index.html")
	})
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandleFunc("/api/", server.Protected(server.ReadSession))

	router.HandleFunc("/api/register", HandleFunc(server.Register))
	router.HandleFunc("/api/login", HandleFunc(server.Login))

	router.HandleFunc("/api/users", server.Protected(server.GetUsers))
	router.HandleFunc("/api/posts", server.Protected(server.GetPosts))
	router.HandleFunc("/api/post", server.Protected(server.Post))
	router.HandleFunc("/api/comment", server.Protected(server.Comment))
	router.HandleFunc("/api/comments", server.Protected(server.GetComments))

	// server.Upgrader = websocket.Upgrader{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// 	CheckOrigin: func(r *http.Request) bool {
	// 		return true
	// 	},
	// }
	// router.HandleFunc("/api/ws", (HandleFunc(server.WS)))

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

func (server *API) Protected(fn HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := server.Sessions.GetSession(r)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, err.Error())
		}
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
		return writeJSON(writer, http.StatusBadRequest, "Invalid Email")
	}

	if _, err = strconv.Atoi(registerReq.Age); err != nil {
		return writeJSON(writer, http.StatusBadRequest, "Age field is invalid")
	}

	var cancel context.CancelFunc
	registerReq.Ctx, cancel = context.WithTimeout(request.Context(), database.TransactionTimeout)
	defer cancel()
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

	if _, err = mail.ParseAddress(loginReq.Email); err != nil {
		return writeJSON(writer, http.StatusBadRequest, "Invalid Email")
	}
	var cancel context.CancelFunc
	loginReq.Ctx, cancel = context.WithTimeout(request.Context(), database.TransactionTimeout)
	defer cancel()
	user, err := server.Storage.LogUser(loginReq)
	if err != nil {
		return writeJSON(writer, http.StatusBadRequest, "Invalid Password")
	}

	session := server.Sessions.NewSession(writer, request)
	session.User = user

	return writeJSON(writer, http.StatusOK, user)
}

func (server *API) GetUsers(writer http.ResponseWriter, request *http.Request) error {
	ctx, cancel := context.WithTimeout(request.Context(), database.TransactionTimeout)
	defer cancel()

	users, err := server.Storage.GetUsers(ctx)
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
	postReq.Ctx, cancel = context.WithTimeout(request.Context(), database.TransactionTimeout)
	defer cancel()

	post, err := server.Storage.CreatePost(postReq)
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, post)
}

func (server *API) GetPosts(writer http.ResponseWriter, request *http.Request) error {
	ctx, cancel := context.WithTimeout(request.Context(), database.TransactionTimeout)
	defer cancel()
	posts, err := server.Storage.GetPosts(ctx)
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, posts)
}

func (server *API) Comment(writer http.ResponseWriter, request *http.Request) error {
	if request.Method != http.MethodPost {
		return writeJSON(writer, http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	commentReq := new(models.CommentRequest)
	err := json.NewDecoder(request.Body).Decode(commentReq)
	if err != nil {
		return err
	}

	if commentReq.Content == "" ||
		commentReq.PostID == "" {

		return writeJSON(writer, http.StatusBadRequest, "Missing Credentials")
	}

	session, err := server.Sessions.GetSession(request)
	if err != nil {
		return err
	}

	commentReq.UserID = session.User.ID
	var cancel context.CancelFunc
	commentReq.Ctx, cancel = context.WithTimeout(request.Context(), database.TransactionTimeout)
	defer cancel()

	comment, err := server.Storage.CreateComment(commentReq)
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, comment)
}

func (server *API) GetComments(writer http.ResponseWriter, request *http.Request) error {
	ctx, cancel := context.WithTimeout(request.Context(), database.TransactionTimeout)
	defer cancel()
	comments, err := server.Storage.GetComments(ctx)
	if err != nil {
		return err
	}

	return writeJSON(writer, http.StatusOK, comments)
}

// func (server *API) WS(writer http.ResponseWriter, request *http.Request) error {
// 	conn, err := server.Upgrader.Upgrade(writer, request, nil)
// 	if err != nil {
// 		return writeJSON(writer, http.StatusInternalServerError, "Error: connecting to the WebSocket.")
// 	}
// 	conn.Close()
// 	return nil
// }
