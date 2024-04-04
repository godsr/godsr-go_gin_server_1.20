package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

type Message struct {
	Text string `json:"text"`
}

func WsGet(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var message Message
		err = json.Unmarshal(p, &message)

		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Received message: %s\n", message.Text)

		var echoMessage Message
		echoMessage.Text = "응답성공"
		response, err := json.Marshal(echoMessage)
		if err != nil {
			log.Println(err)
			return
		}

		// Echo back the received message
		if err := conn.WriteMessage(messageType, response); err != nil {
			log.Println(err)
			return
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
