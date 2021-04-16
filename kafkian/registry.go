package kafkian

import (
	"fmt"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var reg []record

type record struct {
	consumer *kafka.Consumer
	topics   []string
}

func detachTopicFromConsumer(topic string, c *kafka.Consumer) {
	for _, r := range reg {
		if r.consumer != c {
			continue
		}
		for k, v := range r.topics {
			if v == topic {
				r.topics[k] = r.topics[len(r.topics)-1]
				r.topics = r.topics[:len(r.topics)-1]
				fmt.Println("Detaching topic", v)
				if len(r.topics) == 0 {
					unregister(r.consumer)
				}
				return
			}
		}
	}
}

func getConsumerForTopic(topic string) (c *kafka.Consumer) {
	for _, r := range reg {
		for _, t := range r.topics {
			if t == topic {
				fmt.Println("Consumer found")
				c = r.consumer
				return
			}
		}
	}
	return
}

func isRegistered(c *kafka.Consumer) (ok bool) {
	for _, r := range reg {
		if r.consumer == c {
			ok = true
			return
		}
	}
	return
}

func register(c *kafka.Consumer, topics []string) {
	if isRegistered(c) {
		return
	}
	fmt.Println("Registering consumer")

	reg = append(reg, record{c, topics})
}

func unregister(c *kafka.Consumer) {
	for k, v := range reg {
		if v.consumer == c {
			fmt.Println("Unregistering consumer")
			reg[k] = reg[len(reg)-1]
			reg = reg[:len(reg)-1]
			return
		}
	}
}

func init() {
	reg = []record{}
}
