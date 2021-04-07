package kafkian

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jwmwalrus/bnp"
	"github.com/jwmwalrus/felice-franz/base"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	ws *websocket.Conn
}

var WS *websocket.Conn

// HandleWS handles web socket messages
func HandleWS(w http.ResponseWriter, r *http.Request) {
	var err error
	WS, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// ... Use conn to send and receive messages.
	c := connection{ws: WS}
	go c.processMessages()
}

// processMessages reads and processes messages from the websocket
func (c connection) processMessages() {
	defer func() {
		c.ws.Close()
	}()
	// c.ws.SetReadLimit(maxMessageSize)
	// c.ws.SetReadDeadline(time.Now().Add(pongWait))
	// c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	type consumeMsg struct {
		MsgType string   `json:"type"`
		Env     string   `json:"env"`
		Payload []string `json:"payload"`
	}
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		res := consumeMsg{}
		json.Unmarshal(msg, &res)

		switch res.MsgType {
		case "consume":
			env := base.Conf.GetEnvConfig(res.Env)
			topics, _ := env.FindTopicValues(res.Payload)
			err := CreateConsumers(env, topics)
			bnp.LogOnError(err)
		default:
		}
	}
}
