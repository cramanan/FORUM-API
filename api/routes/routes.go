package routes

import (
	"database/sql"
	"errors"
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
}

func RegisterClient(writer http.ResponseWriter, request *http.Request) {

	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}
	var err error

	user := models.User{
		Email:     request.FormValue("register-email"),
		Name:      request.FormValue("register-username"),
		Password:  request.FormValue("register-password"),
		Gender:    request.FormValue("register-gender"),
		FirstName: request.FormValue("register-first-name"),
		LastName:  request.FormValue("register-last-name"),
	}

	if user.Email == "" ||
		user.Name == "" ||
		user.Password == "" ||
		user.Gender == "" ||
		user.FirstName == "" ||
		user.LastName == "" {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("NO CREDENTIALS")
		return
	}

	_, err = mail.ParseAddress(user.Email)
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

	crypt, err := bcrypt.GenerateFromPassword([]byte(user.Password), 11)
	if err != nil {
		log.Println(err)
		return
	}

	user.Password = string(crypt)
	err = database.AddUser(user)
	if err != nil {
		log.Println(err)
		return
	}

	sess := database.CreateSession(writer, request)
	sess.SetID(user.B64)
	sess.SetName(user.Name)
}

func LogClientIn(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("WRONG REQUEST TYPE")
		return
	}

	user := models.User{
		Email:    request.FormValue("login-email"),
		Password: request.FormValue("login-password"),
	}

	if user.Password == "" {
		log.Println("No Password")
		return
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	comp, err := database.GetUserFromMail(user.Email)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(comp.Password), []byte(user.Password))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	user = *comp
	sess := database.CreateSession(writer, request)
	sess.SetID(user.B64)
	sess.SetName(user.Name)
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
	writer.WriteHeader(http.StatusBadRequest)
}

func Logout(writer http.ResponseWriter, request *http.Request) {
}

func WS(writer http.ResponseWriter, request *http.Request) {
	type WSMessage struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("Error")
		return
	}

	go func() {
		conn.WriteJSON(WSMessage{
			Type: "ping",
		})
		for {

		}
	}()
}
