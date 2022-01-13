package kafkian

import (
	"errors"
	"strings"
	"time"

	"github.com/jwmwalrus/felice-n-franz/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// SubscribeConsumer creates a consumer for the given environment topics
func SubscribeConsumer(env base.Environment, topic string) (err error) {
	var c *kafka.Consumer

	if assoc := getConsumerForTopic(topic); assoc != nil {
		return
	}

	cct := time.Now()
	config := cloneConfiguration(env, topic)
	c, err = kafka.NewConsumer(&config)
	if err != nil {
		return
	}

	reg := register(c, env.Name, topic, "")

	go func() {
		defer c.Close()

		c.Subscribe(topic, nil)

	subscribeLoop:
		for {
			select {
			case <-reg.quit:
				log.Info("Consumer is not registered anymore!")
				if c != nil {
					err := c.Unsubscribe()
					onerror.Warn(err)
				}
				break subscribeLoop
			default:
				ev := c.Poll(100)

				switch e := ev.(type) {
				case *kafka.Message:
					if e.Timestamp.Unix() <= cct.Unix() {
						continue
					}
					log.Infof("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
					if e.Headers != nil {
						log.Infof("%% Headers: %v\n", e.Headers)
					}
					sendKafkaMessage(e, "")
				case kafka.Error:
					log.WithFields(log.Fields{"code": e.Code()}).Error(e.String())
					toast := toastMsg{
						ToastType:    toastError,
						Title:        "Consumer Error",
						Message:      e.String(),
						CanBeIgnored: errorCanBeIgnored(e),
					}
					toast.send()
					if e.Code() == kafka.ErrAllBrokersDown {
						break
					}
				case nil:
				default:
					log.Infof("Ignored %v\n", e)
				}
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
	config := cloneConfiguration(env, topic)
	config["group.id"] = topic
	config["enable.auto.commit"] = false
	config["auto.offset.reset"] = "smallest"
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

	reg := register(c, env.Name, topic, "")

	go func() {
		defer c.Close()

		for _, pm := range (*meta).Topics[topic].Partitions {
			tp := kafka.TopicPartition{Topic: &topic, Partition: pm.ID}
			err := c.IncrementalAssign([]kafka.TopicPartition{tp})
			onerror.Log(err)
		}

	assignLoop:
		for {
			select {
			case <-reg.quit:
				log.Info("Consumer is not registered anymore!")
				err := c.Unsubscribe()
				onerror.Warn(err)
				break assignLoop
			default:
				ev := c.Poll(100)

				switch e := ev.(type) {
				case *kafka.Message:
					if e.Timestamp.Unix() <= cct.Unix() {
						continue
					}
					log.Infof("Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
					if e.Headers != nil {
						log.Infof("Headers: %v\n", e.Headers)
					}
					sendKafkaMessage(e, "")
				case kafka.Error:
					log.WithFields(log.Fields{"code": e.Code()}).Error(e.String())
					toast := toastMsg{
						ToastType:    toastError,
						Title:        "Consumer Error",
						Message:      e.String(),
						CanBeIgnored: errorCanBeIgnored(e),
					}
					toast.send()
					if e.Code() == kafka.ErrAllBrokersDown {
						break
					}
				case nil:
				default:
					log.Infof("Ignored %v\n", e)
				}
			}
		}
	}()
	return
}

func cloneConfiguration(env base.Environment, topic string) (out kafka.ConfigMap) {
	out = kafka.ConfigMap{}
	for k, v := range env.Configuration {
		out[k] = v
	}
	if t, err := env.GetTopic(topic); err == nil {
		if t.GroupID != "" {
			out["group.id"] = t.GroupID
		}
	}
	return
}

func errorCanBeIgnored(e kafka.Error) bool {
	return strings.HasPrefix(e.String(), "FindCoordinator") || strings.Contains(e.Code().String(), "Timed out")
}
