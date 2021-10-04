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
	searchID string
	quit     chan struct{}
}

func getConsumerForSearchSearchID(searchID string) (c *kafka.Consumer) {
	for _, r := range reg {
		if r.searchID == searchID {
			log.Info("Consumer found for searchID " + searchID)
			c = r.consumer
			return
		}
	}
	return
}

func getConsumerForTopic(topic string) (c *kafka.Consumer) {
	for _, r := range reg {
		if r.topic == topic && r.searchID == "" {
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

func register(c *kafka.Consumer, env, topic, searchID string) *record {
	if isRegistered(c) {
		return nil
	}
	log.Info("Registering consumer")

	r := record{c, env, topic, searchID, make(chan struct{})}
	reg = append(reg, r)

	return &r
}

func refreshRegistry() {
	log.Info("Refreshing registry")
	var reg0 []record
	copy(reg0, reg)
	for _, r := range reg0 {
		if r.searchID != "" {
			continue
		}
		unregister(r.consumer)
	}

	allSuccess := true
	for _, r := range reg0 {
		if r.searchID != "" {
			continue
		}
		env := base.Conf.GetEnvConfig(r.env)
		if err := SubscribeConsumer(env, r.topic); err != nil {
			allSuccess = false
			log.Error(err)
			toast := toastMsg{
				ToastType: toastError,
				Title:     "Refresh error",
				Message:   err.Error(),
			}
			toast.send()
		}
	}

	if allSuccess {
		toast := toastMsg{
			ToastType: toastInfo,
			Title:     "From Refresh",
			Message:   "All consumers successfully refreshed!",
		}
		toast.send()
	}
}

func resetRegistry() {
	log.Info("Resetting registry")
	for _, v := range reg {
		v.quit <- struct{}{}
	}
	reg = []record{}
}

func unregister(c *kafka.Consumer) {
	for k, v := range reg {
		if v.consumer == c {
			reg[k] = reg[len(reg)-1]
			reg = reg[:len(reg)-1]
			log.Info("Consumer unregistered")
			v.quit <- struct{}{}
			return
		}
	}
}

func init() {
	reg = []record{}
}
