package routes

import (
	"backend/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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
	fmt.Println(request.Method)
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}

	err := request.ParseForm()
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(utils.SERVERMESSAGES[http.StatusBadRequest])
		return
	}

	
}
