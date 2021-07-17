package kafkian

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jwmwalrus/bnp"
	"github.com/jwmwalrus/felice-n-franz/internal/base"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// WS web socket connection
var WS *websocket.Conn

var socketGuard chan struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

type connection struct {
	ws *websocket.Conn
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

			if env.AssignConsumer {
				for _, t := range topics {
					if err := AssignConsumer(env, t); err != nil {
						log.Error(err)
						toast := toastMsg{
							Title:   "Consumer Error",
							Message: err.Error(),
						}
						toast.send()
					}
				}
			} else {
				if err := SubscribeConsumer(env, topics); err != nil {
					log.Error(err)
					toast := toastMsg{
						Title:   "Consumer Error",
						Message: err.Error(),
					}
					toast.send()
				}
			}

		case "produce":
			env := base.Conf.GetEnvConfig(res.Env)
			msg := kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &res.Topic},
				Value:          []byte(res.Payload[0]),
			}
			if res.Key != "" {
				msg.Key = []byte(res.Key)
			}
			if len(res.Headers) > 0 {
				msg.Headers = headerB2K(res.Headers)
			}
			if err := ProduceMessage(env, &msg); err != nil {
				log.Error(err)
				toast := toastMsg{
					Title:   "Producer Error",
					Message: err.Error(),
				}
				toast.send()
			}

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

type receivedMsg struct {
	MsgType string        `json:"type"`
	Env     string        `json:"env"`
	Topic   string        `json:"topic"`
	Key     string        `json:"key"`
	Payload []string      `json:"payload"`
	Headers []base.Header `json:"headers"`
}

type toastMsg struct {
	ToastType    string       `json:"toastType"`
	Title        string       `json:"title"`
	Message      string       `json:"message"`
	Topic        string       `json:"topic"`
	Partition    int32        `json:"partition"`
	Offset       kafka.Offset `json:"offset"`
	CanBeIgnored bool         `json:"canBeIgnored"`
	SentAt       int64        `json:"sentAt"`
}

func (t toastMsg) send() {
	t.SentAt = time.Now().Unix()
	payload, err := json.Marshal(t)
	if err != nil {
		log.Error(err)
		return
	}

	socketGuard <- struct{}{}
	go writeToSocket(payload)
	return
}

func sendKafkaMessage(m *kafka.Message) {
	headers := headerK2B(m.Headers)

	flat := struct {
		Topic               string              `json:"topic"`
		Partition           int32               `json:"partition"`
		Offset              kafka.Offset        `json:"offset"`
		Value               string              `json:"value"`
		Key                 string              `json:"key"`
		Timestamp           time.Time           `json:"timestamp"`
		TimestampType       kafka.TimestampType `json:"timestampType"`
		TimestampTypeString string              `json:"timestampTypeString"`
		Headers             []base.Header       `json:"headers"`
	}{
		*m.TopicPartition.Topic,
		m.TopicPartition.Partition,
		m.TopicPartition.Offset,
		string(m.Value),
		string(m.Key),
		m.Timestamp,
		m.TimestampType,
		m.TimestampType.String(),
		headers,
	}
	payload, err := json.Marshal(flat)
	if err != nil {
		log.Error(err)
		return
	}

	socketGuard <- struct{}{}
	go writeToSocket(payload)
	return
}

func writeToSocket(payload []byte) {
	err := WS.WriteMessage(websocket.TextMessage, payload)
	bnp.LogOnError(err)
	<-socketGuard
}

func headerK2B(kh []kafka.Header) (h []base.Header) {
	for _, v := range kh {
		h = append(h, base.Header{v.Key, string(v.Value)})
	}
	return
}

func headerB2K(h []base.Header) (kh []kafka.Header) {
	for _, v := range h {
		kh = append(kh, kafka.Header{v.Key, []byte(v.Value)})
	}
	return
}

func init() {
	socketGuard = make(chan struct{}, 1)
}
