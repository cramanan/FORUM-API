package routes

import (
	"log"
	"net/http"
	"net/mail"
	"real-time-forum/api/database"
	"real-time-forum/api/models"
	"real-time-forum/api/utils"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Root(writer http.ResponseWriter, request *http.Request) error {
	sess := contextSession(request)
	if sess == nil {
		return writeJSON(writer, http.StatusServiceUnavailable, nil)
	}
	return writeJSON(writer, http.StatusOK, sess.GetUser())
}

func Register(writer http.ResponseWriter, request *http.Request) error {
	if request.Method != http.MethodPost {
		return writeJSON(writer, http.StatusMethodNotAllowed, nil)
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
		return writeJSON(writer, http.StatusBadRequest, nil)
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return writeJSON(writer, http.StatusBadRequest, nil)
	}

	user.Age, err = strconv.Atoi(request.FormValue("register-age"))
	if err != nil {
		return writeJSON(writer, http.StatusBadRequest, nil)
	}

	user.B64 = utils.GenerateBase64ID(5)
	crypt, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return writeJSON(writer, http.StatusInternalServerError, nil)
	}
	user.SetPassword(crypt)

	err = database.AddUser(user)
	if err != nil {
		return writeJSON(writer, http.StatusInternalServerError, nil)
	}

	sess := database.NewSession(writer, request)
	sess.SetUser(user)
	return writeJSON(writer, http.StatusOK, nil)
}

func Login(writer http.ResponseWriter, request *http.Request) error {
	if request.Method != http.MethodPost {
		return writeJSON(writer, http.StatusMethodNotAllowed, nil)
	}

	email := request.FormValue("login-email")
	password := []byte(request.FormValue("login-password"))
	parsedMail, err := mail.ParseAddress(email)
	if err != nil {
		return writeJSON(writer, http.StatusBadRequest, nil)
	}

	comp, err := database.GetUserFromMail(parsedMail)
	if err != nil {
		log.Println(err)
		return writeJSON(writer, http.StatusBadRequest, nil)
	}

	err = bcrypt.CompareHashAndPassword(comp.GetPassword(), password)
	if err != nil {
		log.Println("INVALID PASSWORD")
		return writeJSON(writer, http.StatusInternalServerError, nil)
	}

	sess := database.NewSession(writer, request)
	sess.SetUser(comp)
	return writeJSON(writer, http.StatusOK, nil)
}

func Post(writer http.ResponseWriter, request *http.Request) error {
	if request.Method != http.MethodPost {
		return writeJSON(writer, http.StatusMethodNotAllowed, nil)
	}

	sess := contextSession(request)
	if sess == nil {
		log.Println("NO SESSION")
		return writeJSON(writer, http.StatusServiceUnavailable, nil)
	}

	p := models.Post{
		UserID:   sess.GetUser().B64,
		Username: sess.GetUser().Name,
		Content:  request.FormValue("post-content"),
		Date:     time.Now(),
	}

	err := database.CreatePost(p)
	if err != nil {
		log.Println(err)
		return writeJSON(writer, http.StatusInternalServerError, nil)
	}
	return writeJSON(writer, http.StatusOK, nil)
}

func GetPosts(writer http.ResponseWriter, request *http.Request) error {
	posts, err := database.GetAllPosts()
	if err != nil {
		return writeJSON(writer, http.StatusServiceUnavailable, nil)
	}
	return writeJSON(writer, http.StatusOK, posts)
}

func GetUsers(writer http.ResponseWriter, request *http.Request) error {
	users, err := database.GetAllUsers()
	if err != nil {
		return writeJSON(writer, http.StatusServiceUnavailable, nil)
	}
	return writeJSON(writer, http.StatusOK, users)
}

func Logout(writer http.ResponseWriter, request *http.Request) error {
	sess := contextSession(request)
	if sess == nil {
		log.Println("NO SESSION")
		return writeJSON(writer, http.StatusServiceUnavailable, nil)
	}
	sess.End()
	return writeJSON(writer, http.StatusOK, nil)

}

// var (
// 	upgrader = websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}
// )

// func WS(writer http.ResponseWriter, request *http.Request) error {
// 	type WSMessage struct {
// 		Type string      `json:"type"`
// 		Data interface{} `json:"data"`
// 	}

// 	conn, err := upgrader.Upgrade(writer, request, nil)
// 	if err != nil {
// 		log.Println("Error Upgrading protocol")
// 		return writeJSON(writer, http.StatusServiceUnavailable, nil)

// 	}

// 	conn.WriteJSON(WSMessage{
// 		Type: "ping",
// 		Data: "Hello",
// 	})
// 	go func() {
// 		for {

// 		}
// 	}()
// 	return writeJSON(writer, http.StatusOK, nil)
// }
