package routes

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"fmt"
	"log"
	"net/http"
	"net/mail"

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
	writer.Header().Add("Content-Type", "application/json")
	_, err := Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {

	}

}

func RegisterClient(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		fmt.Println("Wrong Request")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}

	err := request.ParseForm()
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}

	client := models.Client{
		Username: request.FormValue("username"),
		Password: request.FormValue("password"),
	}

	if client.Username == "" || client.Password == "" {
		fmt.Println("Empty credentials")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}

	client.Email, err = mail.ParseAddress(request.FormValue("email"))
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}

	client.Uuid, err = uuid.NewV4()
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(utils.SERVERMESSAGES[http.StatusInternalServerError])
		return
	}

	crypt, err := bcrypt.GenerateFromPassword([]byte(client.Password), 11)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(utils.SERVERMESSAGES[http.StatusInternalServerError])
		return
	}

	client.Password = string(crypt)

	err = database.AddClient(client)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(utils.SERVERMESSAGES[http.StatusInternalServerError])
		return
	}
}

func LogClientIn(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		fmt.Println("Wrong Request")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}

	err := request.ParseForm()
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}

	client := models.Client{
		Username: request.FormValue("username"),
		Password: request.FormValue("password"),
	}

	if client.Username == "" || client.Password == "" {
		fmt.Println("Empty credentials")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}

	client.Email, err = mail.ParseAddress(request.FormValue("email"))
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}
}
