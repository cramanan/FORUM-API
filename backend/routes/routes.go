package routes

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strconv"

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

	raw, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		resp.StatusCode = http.StatusInternalServerError
		resp.Message = "Something Went Wrong :/ Try again later."
		SendResponse(writer, resp)
		return
	}

	user.Uuid = raw.String()

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
	sess.Set("uuid", user.Uuid)
	sess.Set("username", user.Username)
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
	sess.Set("uuid", user.Uuid)
	sess.Set("username", user.Username)
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

func Logout(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	writer.Header().Set("Access-Control-Allow-Credentials", "true")

}
