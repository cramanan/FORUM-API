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
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	resp := Response{
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
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	resp := Response{
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

	client := models.Client{
		Username: request.FormValue("register-username"),
		Password: request.FormValue("register-password"),
		Gender:   request.FormValue("register-gender"),
		//Add Age
		FirstName: request.FormValue("register-first-name"),
		LastName:  request.FormValue("register-last-name"),
	}

	if client.Username == "" ||
		client.Password == "" ||
		client.Gender == "" ||
		client.FirstName == "" ||
		client.LastName == "" {

		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Empty Credentials"
		SendResponse(writer, resp)
		return
	}

	client.Email, err = mail.ParseAddress(request.FormValue("register-email"))
	if err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Invalid Email format"
		SendResponse(writer, resp)
		return
	}

	client.Uuid, err = uuid.NewV4()
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Something Went Wrong :/ Try again later."
		SendResponse(writer, resp)
		return
	}

	crypt, err := bcrypt.GenerateFromPassword([]byte(client.Password), 11)
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Something Went Wrong :/ Try again later."
		SendResponse(writer, resp)
		return
	}

	client.Password = string(crypt)
	err = database.AddClient(client)
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
		SendResponse(writer, resp)
		return
	}

	sess := database.CreateSession(writer, request)
	sess.Set("uuid", client.Uuid.String())
	sess.Set("username", client.Username)
	SendResponse(writer, resp)
}

func LogClientIn(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	resp := Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}
	if request.Method != http.MethodPost {
		resp.StatusCode = http.StatusBadRequest
		resp.Message = "Bad Request"
		SendResponse(writer, resp)
		return
	}

	client := models.Client{
		Password: request.FormValue("login-password"),
	}

	if client.Password == "" {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Empty credentials"
		SendResponse(writer, resp)
		return
	}
	var err error
	client.Email, err = mail.ParseAddress(request.FormValue("login-email"))
	if err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Invalid mail format"
		SendResponse(writer, resp)
		return
	}

	comp, err := database.GetClientFromMail(client.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp.StatusCode = http.StatusUnauthorized
			resp.Message = "Invalid password or username"
			SendResponse(writer, resp)
			return
		}
		resp.StatusCode = http.StatusInternalServerError
		SendResponse(writer, resp)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(comp.Password), []byte(client.Password))
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Invalid password or username"
		SendResponse(writer, resp)
		return
	}

	client = *comp
	sess := database.CreateSession(writer, request)
	sess.Set("uuid", client.Uuid.String())
	sess.Set("username", client.Username)
	SendResponse(writer, resp)
}

func Post(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	resp := Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
	}

	if request.Method != http.MethodPost {
		resp.StatusCode = http.StatusBadRequest
		resp.Message = "Bad Request"
		SendResponse(writer, resp)
		return
	}

	sess, err := database.GetSession(writer, request)
	if err != nil {
		resp.StatusCode = http.StatusUnauthorized
		resp.Message = "Unauthorized"
		SendResponse(writer, resp)
		return
	}

	p := models.Post{}
	assertion, ok := sess.Get("uuid")
	if !ok {
		log.Println("UUID not found")
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
		SendResponse(writer, resp)
		return
	}

	p.UserID, ok = assertion.(string)
	if !ok {
		log.Println("UUID not a string")
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
		SendResponse(writer, resp)
		return
	}

	p.Content = request.FormValue("post-content")
	if p.Content == "" {
		log.Println("Empty content")
		resp.StatusCode = http.StatusBadRequest
		resp.Message = "Bad Request"
		SendResponse(writer, resp)
		return
	}

	err = database.CreatePost(p)
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Internal Server Error"
	}
	SendResponse(writer, resp)
}

func GetPosts(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	writer.Header().Set("Access-Control-Allow-Credentials", "true")
	log.Println("Server Reached")
	posts, err := database.GetAllPosts()
	if err != nil {
		log.Println(err)
		return
	}

	resp := Response{
		StatusCode: http.StatusOK,
		Message:    "OK",
		Data:       posts,
	}
	SendResponse(writer, resp)
}
