package kafkian

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jwmwalrus/bnp"
	"github.com/jwmwalrus/felice-franz/base"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// CreateConsumer creates a consumer for the given environment topics
func CreateConsumer(env base.Environment, topics []string) (err error) {
	var c *kafka.Consumer

	c, err = kafka.NewConsumer(&env.Configuration)

	if err != nil {
		return
	}

	go func() {
		defer c.Close()

		c.SubscribeTopics(topics, nil)

		for {
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
				payload, err := getPayLoadFromMessage(e)
				if err != nil {
					log.Error(err)
					continue
				}
				err = WS.WriteMessage(websocket.TextMessage, payload)
				bnp.LogOnError(err)
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				log.WithFields(log.Fields{"code": e.Code()}).Error(e.String())
				if e.Code() == kafka.ErrAllBrokersDown {
					break
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}()
	return
}

func getPayLoadFromMessage(m *kafka.Message) (payload []byte, err error) {
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
		Partition     int32               `json:"partitio"`
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
	payload, err = json.Marshal(flat)
	return
}
