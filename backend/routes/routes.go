package routes

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/mail"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Root(writer http.ResponseWriter, request *http.Request) {

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

	models.Session(writer, request)
	err = models.SendResponse(writer, resp)
	log.Println(err)
}

func LogClientIn(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
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

	client = *comp
	models.Session(writer, request)
	models.SendResponse(writer, resp)
}
