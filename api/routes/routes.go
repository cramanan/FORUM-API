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
	resp := &Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
	sess, err := database.GetSession(writer, request)
	if err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Unauthorized"
	} else {
		resp.Data = sess.Values()
	}
	SendResponse(writer, resp)
}

func RegisterClient(writer http.ResponseWriter, request *http.Request) {
	resp := &Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
	if request.Method != http.MethodPost {
		resp.StatusCode = http.StatusBadRequest
		resp.Message = "Bad Request"
		SendResponse(writer, resp)
		return
	}
	var err error

	user := models.User{
		Email:     request.FormValue("register-email"),
		Username:  request.FormValue("register-username"),
		Password:  request.FormValue("register-password"),
		Gender:    request.FormValue("register-gender"),
		FirstName: request.FormValue("register-first-name"),
		LastName:  request.FormValue("register-last-name"),
	}

	if user.Email == "" ||
		user.Username == "" ||
		user.Password == "" ||
		user.Gender == "" ||
		user.FirstName == "" ||
		user.LastName == "" {

		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Empty Credentials"
		SendResponse(writer, resp)
		return
	}

	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Invalid Email format"
		SendResponse(writer, resp)
		return
	}

	user.Age, err = strconv.Atoi(request.FormValue("register-age"))
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Invalid Age format"
		SendResponse(writer, resp)
		return
	}

	user.B64 = utils.GenerateBase64ID(5)

	crypt, err := bcrypt.GenerateFromPassword([]byte(user.Password), 11)
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Something Went Wrong :/ Try again later."
		SendResponse(writer, resp)
		return
	}

	user.Password = string(crypt)
	err = database.AddClient(user)
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
		SendResponse(writer, resp)
		return
	}

	sess := database.CreateSession(writer, request)
	sess.Set("b64", user.B64)
	sess.Set("username", user.Username)
	SendResponse(writer, resp)
}

func LogClientIn(writer http.ResponseWriter, request *http.Request) {
	resp := &Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
	if request.Method != http.MethodPost {
		resp.StatusCode = http.StatusBadRequest
		resp.Message = "Bad Request"
		SendResponse(writer, resp)
		return
	}

	user := models.User{
		Email:    request.FormValue("login-email"),
		Password: request.FormValue("login-password"),
	}

	if user.Password == "" {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Empty credentials"
		SendResponse(writer, resp)
		return
	}

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Invalid mail format"
		SendResponse(writer, resp)
		return
	}

	comp, err := database.GetClientFromMail(user.Email)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			resp.StatusCode = http.StatusUnauthorized
			resp.Message = "Invalid password or email"
			SendResponse(writer, resp)
			return
		}
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
		SendResponse(writer, resp)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(comp.Password), []byte(user.Password))
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Invalid password or email"
		SendResponse(writer, resp)
		return
	}

	user = *comp
	sess := database.CreateSession(writer, request)
	sess.Set("b64", user.B64)
	sess.Set("username", user.Username)
	SendResponse(writer, resp)
}

func Post(writer http.ResponseWriter, request *http.Request) {
	resp := &Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}

	if request.Method != http.MethodPost {
		resp.StatusCode = http.StatusBadRequest
		resp.Message = "Bad Request"
		SendResponse(writer, resp)
		return
	}

	session, ok := request.Context().Value(middleware.ContextSessionKey).(*database.Session)
	if !ok {
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
		SendResponse(writer, resp)
		return
	}

	log.Println(session)

	SendResponse(writer, resp)
}

func GetPosts(writer http.ResponseWriter, request *http.Request) {
	// Check session
	posts, err := database.GetAllPosts()
	if err != nil {
		log.Println(err)
		return
	}

	resp := &Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data:       posts,
	}
	SendResponse(writer, resp)
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
		SendResponse(writer, nil)
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
