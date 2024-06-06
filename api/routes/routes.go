package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"real-time-forum/api/database"
	"real-time-forum/api/models"
	"real-time-forum/api/routes/middleware"
	"real-time-forum/api/utils"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func Root(writer http.ResponseWriter, request *http.Request) {
	sess, ok := request.Context().Value(middleware.ContextSessionKey).(*database.Session)
	if !ok {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		json.NewEncoder(writer).Encode(sess.GetName())
	}
}

func RegisterUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{
		Email:     request.FormValue("register-email"),
		Name:      request.FormValue("register-username"),
		Gender:    request.FormValue("register-gender"),
		FirstName: request.FormValue("register-first-name"),
		LastName:  request.FormValue("register-last-name"),
	}

	password := request.FormValue("register-password")
	if user.Email == "" ||
		user.Name == "" ||
		password == "" ||
		user.Gender == "" ||
		user.FirstName == "" ||
		user.LastName == "" {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("NO CREDENTIALS")
		return
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("INVALID EMAIL")
		return
	}

	user.Age, err = strconv.Atoi(request.FormValue("register-age"))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	user.B64 = utils.GenerateBase64ID(5)

	crypt, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		log.Println(err)
		return
	}

	user.SetPassword(crypt)
	err = database.AddUser(user)
	if err != nil {
		log.Println(err)
		return
	}

	sess := database.NewSession(writer, request)
	sess.SetID(user.B64)
	sess.SetName(user.Name)
}

func LogUserIn(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("WRONG REQUEST TYPE")
		return
	}

	name := request.FormValue("login-name")
	email := request.FormValue("login-email")

	password := request.FormValue("login-password")
	if password == "" {
		log.Println("No Password")
		return
	}

	parsedMail, err := mail.ParseAddress(email)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("INVALID EMAIL")
		return
	}

	comp, b64, err := database.GetPasswordAndIDFromMail(*parsedMail)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword(comp, []byte(password))
	if err != nil {
		log.Println("INVALID PASSWORD")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	sess := database.NewSession(writer, request)
	sess.SetID(b64)
	sess.SetName(name)
}

func Post(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("WRONG REQUEST TYPE")
		return
	}

	sess, ok := request.Context().Value(middleware.ContextSessionKey).(*database.Session)
	if !ok {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println("NO SESSION")
		return
	}

	p := models.Post{
		UserID:   sess.GetID(),
		Username: sess.GetName(),
		Content:  request.FormValue("post-content"),
		Date:     time.Now(),
	}

	err := database.CreatePost(p)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func GetPosts(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	posts, err := database.GetAllPosts()
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(writer).Encode(posts)
}

func GetUsers(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	users, err := database.GetAllUsers()
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(writer).Encode(users)
}

func Logout(writer http.ResponseWriter, request *http.Request) {
	sess, ok := request.Context().Value(middleware.ContextSessionKey).(*database.Session)
	if !ok {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println("NO SESSION")
		return
	}
	sess.End()
}

func WS(writer http.ResponseWriter, request *http.Request) {
	type WSMessage struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("Error Upgrading protocol")
		return
	}

	conn.WriteJSON(WSMessage{
		Type: "ping",
		Data: "Hello",
	})
	go func() {
		for {

		}
	}()
}
