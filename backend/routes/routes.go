package routes

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func BasicUpgrade(writer http.ResponseWriter, request *http.Request) {
	_, err := Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {

	}
}

func RegisterClient(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	resp := models.Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
	if request.Method != http.MethodPost {
		resp.StatusCode = http.StatusBadRequest
		resp.Message = "Bad Request"
		models.SendResponse(writer, resp)
		return
	}
	var err error

	client := models.Client{
		Username: request.FormValue("username"),
		Password: request.FormValue("password"),
	}

	if client.Username == "" || client.Password == "" {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Empty Credentials"
		models.SendResponse(writer, resp)
		return
	}

	client.Email, err = mail.ParseAddress(request.FormValue("email"))

	if err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Invalid Email format"
		models.SendResponse(writer, resp)
		return
	}

	client.Uuid, err = uuid.NewV4()
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Something Went Wrong :/ Try again later."
		models.SendResponse(writer, resp)
		return
	}

	crypt, err := bcrypt.GenerateFromPassword([]byte(client.Password), 11)
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Something Went Wrong :/ Try again later."
		models.SendResponse(writer, resp)
		return
	}

	client.Password = string(crypt)

	err = database.AddClient(client)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: clients.email" {
			resp.StatusCode = http.StatusConflict
			resp.Message = "Email adress already taken"
			models.SendResponse(writer, resp)
			return
		}
		resp.StatusCode = http.StatusInternalServerError
		models.SendResponse(writer, resp)
		return
	}

}

func LogClientIn(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	resp := models.Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
	if request.Method != http.MethodPost {
		resp.StatusCode = http.StatusBadRequest
		resp.Message = "Bad Request"
		models.SendResponse(writer, resp)
		return
	}

	client := models.Client{
		Password: request.FormValue("password"),
	}

	if client.Password == "" {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Empty credentials"
		models.SendResponse(writer, resp)
		return
	}
	var err error
	client.Email, err = mail.ParseAddress(request.FormValue("email"))
	if err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Invalid mail format"
		models.SendResponse(writer, resp)
		return
	}

	comp, err := database.GetClientFromMail(client.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp.StatusCode = http.StatusUnauthorized
			resp.Message = "Invalid password or username"
			models.SendResponse(writer, resp)
			return
		}
		resp.StatusCode = http.StatusInternalServerError
		models.SendResponse(writer, resp)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(comp.Password), []byte(client.Password))
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Invalid password or username"
		models.SendResponse(writer, resp)
		return
	}

	_, ok := utils.Store.Sessions[client.Uuid.String()]
	if ok {
		resp.StatusCode = http.StatusConflict
		resp.Message = "User is already active"
		models.SendResponse(writer, resp)
		return
	}

	session := models.NewSession()
	utils.Store.Sessions[client.Uuid.String()] = session
	cookie := http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	http.SetCookie(writer, &cookie)
	models.SendResponse(writer, resp)
}
