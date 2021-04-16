package kafkian

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jwmwalrus/bnp"
	"github.com/jwmwalrus/felice-franz/base"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// WS web socket connection
var WS *websocket.Conn

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	ws *websocket.Conn
}

type receivedMsg struct {
	MsgType string   `json:"type"`
	Env     string   `json:"env"`
	Payload []string `json:"payload"`
}

type toastMsg struct {
	ToastType string       `json:"toastType"`
	Title     string       `json:"title"`
	Message   string       `json:"message"`
	Topic     string       `json:"topic"`
	Partition int32        `json:"partition"`
	Offset    kafka.Offset `json:"offset"`
}

// HandleWS handles web socket messages
func HandleWS(w http.ResponseWriter, r *http.Request) {
	var err error
	WS, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := connection{ws: WS}
	go c.processMessages()
}

func (c connection) processMessages() {
	defer func() {
		c.ws.Close()
	}()
	// c.ws.SetReadLimit(maxMessageSize)
	// c.ws.SetReadDeadline(time.Now().Add(pongWait))
	// c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		res := receivedMsg{}
		json.Unmarshal(msg, &res)

		switch res.MsgType {
		case "consume":
			env := base.Conf.GetEnvConfig(res.Env)
			topics, _ := env.FindTopicValues(res.Payload)
			err := CreateConsumer(env, topics)
			bnp.LogOnError(err)

		case "produce":
			env := base.Conf.GetEnvConfig(res.Env)
			topic := ""
			payload := []byte{}
			if len(res.Payload) > 1 {
				topic = res.Payload[0]
				payload = []byte(res.Payload[1])
			}
			err := ProduceMessage(env, topic, payload)
			bnp.LogOnError(err)

		case "unsubscribe":
			for _, t := range res.Payload {
				if c := getConsumerForTopic(t); c != nil {
					detachTopicFromConsumer(t, c)
				}
			}
		default:
		}
	}
}

func sendKafkaMessage(m *kafka.Message) {
	type ht struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	headers := []ht{}
	for _, h := range m.Headers {
		headers = append(headers, ht{h.Key, string(h.Value)})
	}

	flat := struct {
		Topic         string              `json:"topic"`
		Partition     int32               `json:"partition"`
		Offset        kafka.Offset        `json:"offset"`
		Value         string              `json:"value"`
		Key           string              `json:"key"`
		Timestamp     time.Time           `json:"timestamp"`
		TimestampType kafka.TimestampType `json:"timestampType"`
		Headers       []ht                `json:"headers"`
	}{
		*m.TopicPartition.Topic,
		m.TopicPartition.Partition,
		m.TopicPartition.Offset,
		string(m.Value),
		string(m.Key),
		m.Timestamp,
		m.TimestampType,
		headers,
	}
	payload, err := json.Marshal(flat)
	if err != nil {
		log.Error(err)
		return
	}
	err = WS.WriteMessage(websocket.TextMessage, payload)
	bnp.LogOnError(err)
	return
}

func sendToast(t toastMsg) {
	payload, err := json.Marshal(t)
	if err != nil {
		log.Error(err)
		return
	}
	err = WS.WriteMessage(websocket.TextMessage, payload)
	bnp.LogOnError(err)
}
