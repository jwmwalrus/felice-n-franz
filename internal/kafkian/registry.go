package kafkian

import (
	"github.com/jwmwalrus/felice-n-franz/internal/base"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var reg []record

type record struct {
	consumer *kafka.Consumer
	env      string
	topic    string
}

func getConsumerForTopic(topic string) (c *kafka.Consumer) {
	for _, r := range reg {
		if r.topic == topic {
			log.Info("Consumer found for topic " + topic)
			c = r.consumer
			return
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

func register(c *kafka.Consumer, env string, topic string) {
	if isRegistered(c) {
		return
	}
	log.Info("Registering consumer")

	reg = append(reg, record{c, env, topic})
}

func refreshRegistry() {
	log.Info("Refreshing registry")
	var reg0 []record
	copy(reg0, reg)
	for _, r := range reg0 {
		unregister(r.consumer)
	}

	for _, r := range reg0 {
		env := base.Conf.GetEnvConfig(r.env)
		if err := SubscribeConsumer(env, r.topic); err != nil {
			log.Error(err)
			toast := toastMsg{
				Title:   "Consumer Error",
				Message: err.Error(),
			}
			toast.send()
		}
	}
}

func resetRegistry() {
	log.Info("Resetting registry")
	reg = []record{}
}

func unregister(c *kafka.Consumer) {
	for k, v := range reg {
		if v.consumer == c {
			reg[k] = reg[len(reg)-1]
			reg = reg[:len(reg)-1]
			log.Info("Consumer unregistered")
			return
		}
	}
}

func init() {
	reg = []record{}
}
