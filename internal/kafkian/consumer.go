package kafkian

import (
	"github.com/jwmwalrus/bnp"
	"github.com/jwmwalrus/felice-n-franz/internal/base"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// CreateConsumer creates a consumer for the given environment topics
func CreateConsumer(env base.Environment, topics []string) (err error) {
	var c *kafka.Consumer

	uncovered := []string{}
	for _, t := range topics {
		if assoc := getConsumerForTopic(t); assoc != nil {
			continue
		}
		uncovered = append(uncovered, t)
	}

	if len(uncovered) == 0 {
		return
	}

	c, err = kafka.NewConsumer(&env.Configuration)
	if err != nil {
		return
	}

	register(c, uncovered)

	go func() {
		defer c.Close()

		c.SubscribeTopics(uncovered, nil)

		for {
			ev := c.Poll(100)
			if !isRegistered(c) {
				log.Info("Unsubscribing consumer")
				err := c.Unsubscribe()
				bnp.WarnOnError(err)
				break
			}

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				log.Infof("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					log.Infof("%% Headers: %v\n", e.Headers)
				}
				sendKafkaMessage(e)
			case kafka.Error:
				log.WithFields(log.Fields{"code": e.Code()}).Error(e.String())
				sendToast(toastMsg{ToastType: "error", Title: "Consumer Error " + e.Code().String(), Message: e.String()})
				if e.Code() == kafka.ErrAllBrokersDown {
					break
				}
			default:
				log.Infof("Ignored %v\n", e)
			}
		}
	}()
	return
}
