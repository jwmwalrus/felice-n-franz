package kafkian

import (
	"errors"
	"strings"
	"time"

	"github.com/jwmwalrus/bnp"
	"github.com/jwmwalrus/felice-n-franz/internal/base"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// SubscribeConsumer creates a consumer for the given environment topics
func SubscribeConsumer(env base.Environment, topics []string) (err error) {
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

	cct := time.Now()
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
				log.Info("Consumer is not registered anymore!")
				err := c.Unsubscribe()
				bnp.WarnOnError(err)
				break
			}

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				if e.Timestamp.Unix() <= cct.Unix() {
					continue
				}
				log.Infof("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					log.Infof("%% Headers: %v\n", e.Headers)
				}
				sendKafkaMessage(e)
			case kafka.Error:
				log.WithFields(log.Fields{"code": e.Code()}).Error(e.String())
				if !errorCanBeIgnored(e) {
					sendError(errors.New(e.String()), "Consumer Error")
				}
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

// AssignConsumer creates a consumer for the given environment topic
func AssignConsumer(env base.Environment, topic string) (err error) {
	var c *kafka.Consumer

	if assoc := getConsumerForTopic(topic); assoc != nil {
		return
	}

	cct := time.Now()
	config := cloneConfiguration(env.Configuration)
	config["group.id"] = topic
	c, err = kafka.NewConsumer(&config)
	if err != nil {
		return
	}

	var meta *kafka.Metadata
	meta, err = c.GetMetadata(&topic, false, 10000)
	if err != nil {
		return
	}

	if len((*meta).Topics[topic].Partitions) == 0 {
		err = errors.New("There are no available partitions for topic " + topic)
		return
	}

	register(c, []string{topic})

	go func() {
		defer c.Close()

		for _, pm := range (*meta).Topics[topic].Partitions {
			tp := kafka.TopicPartition{Topic: &topic, Partition: pm.ID}
			err := c.IncrementalAssign([]kafka.TopicPartition{tp})
			bnp.LogOnError(err)
		}

		for {
			ev := c.Poll(100)
			if !isRegistered(c) {
				log.Info("Consumer is not registered anymore!")
				err := c.Unsubscribe()
				bnp.WarnOnError(err)
				break
			}

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				if e.Timestamp.Unix() <= cct.Unix() {
					continue
				}
				log.Infof("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					log.Infof("%% Headers: %v\n", e.Headers)
				}
				sendKafkaMessage(e)
			case kafka.Error:
				log.WithFields(log.Fields{"code": e.Code()}).Error(e.String())
				if !errorCanBeIgnored(e) {
					sendError(errors.New(e.String()), "Consumer Error")
				}
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

func cloneConfiguration(in kafka.ConfigMap) (out kafka.ConfigMap) {
	out = kafka.ConfigMap{}
	for k, v := range in {
		out[k] = v
	}
	return
}

func errorCanBeIgnored(e kafka.Error) bool {
	return strings.HasPrefix(e.String(), "FindCoordinator") || strings.Contains(e.Code().String(), "Timed out")
}
