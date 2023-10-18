package kafkian

import (
	"log/slog"

	"github.com/jwmwalrus/felice-n-franz/internal/base"
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
			slog.Info("Consumer found", "searchID", searchID)
			c = r.consumer
			return
		}
	}
	return
}

func getConsumerForTopic(topic string) (c *kafka.Consumer) {
	for _, r := range reg {
		if r.topic == topic && r.searchID == "" {
			slog.Info("Consumer found", "topic", topic)
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
	slog.Info("Registering consumer")

	r := record{c, env, topic, searchID, make(chan struct{})}
	reg = append(reg, r)

	return &r
}

func refreshRegistry() {
	slog.Info("Refreshing registry")
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
			slog.With(
				"env", env,
				"topic", r.topic,
				"error", err,
			).Error("Failed to subscribe consumer")
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
	slog.Info("Resetting registry")
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
			slog.Info("Consumer unregistered")
			v.quit <- struct{}{}
			return
		}
	}
}

func init() {
	reg = []record{}
}
