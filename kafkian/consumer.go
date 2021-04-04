package kafkian

import (
	"fmt"
	"os"

	"github.com/jwmwalrus/felice-franz/base"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// CreateConsumers creates consumers for the given environment topics
func CreateConsumers(env base.Environment, topics []string) (err error) {
	var c *kafka.Consumer

	c, err = kafka.NewConsumer(&env.Configuration)

	if err != nil {
		return
	}

	go func() {
		defer c.Close()

		c.SubscribeTopics(topics, nil)

		for {

			msg, err := c.ReadMessage(-1)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			} else {
				// The client will automatically try to recover from all errors.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
	}()
	return
}

// CreateConsumerPoll creates consumers for the given environment topics
func CreateConsumerPoll(env base.Environment, topics []string) (err error) {
	var c *kafka.Consumer

	// cm := getConfigMap(env)
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
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
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
