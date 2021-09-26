package kafkian

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/bnp/slice"
	"github.com/jwmwalrus/felice-n-franz/internal/base"
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

	register(c, env.Name, topic)

	go func() {
		defer c.Close()

		c.Subscribe(topic, nil)

		for {
			ev := c.Poll(100)
			if !isRegistered(c) {
				log.Info("Consumer is not registered anymore!")
				if c != nil {
					err := c.Unsubscribe()
					onerror.Warn(err)
				}
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
				sendKafkaMessage(e, "")
			case kafka.Error:
				log.WithFields(log.Fields{"code": e.Code()}).Error(e.String())
				toast := toastMsg{
					Title:        "Consumer Error",
					Message:      e.String(),
					ToastType:    toastError,
					CanBeIgnored: errorCanBeIgnored(e),
				}
				toast.send()
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

	register(c, env.Name, topic)

	go func() {
		defer c.Close()

		for _, pm := range (*meta).Topics[topic].Partitions {
			tp := kafka.TopicPartition{Topic: &topic, Partition: pm.ID}
			err := c.IncrementalAssign([]kafka.TopicPartition{tp})
			onerror.Log(err)
		}

		for {
			ev := c.Poll(100)
			if !isRegistered(c) {
				log.Info("Consumer is not registered anymore!")
				err := c.Unsubscribe()
				onerror.Warn(err)
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
				log.Infof("Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					log.Infof("Headers: %v\n", e.Headers)
				}
				sendKafkaMessage(e, "")
			case kafka.Error:
				log.WithFields(log.Fields{"code": e.Code()}).Error(e.String())
				toast := toastMsg{
					Title:        "Consumer Error",
					Message:      e.String(),
					ToastType:    toastError,
					CanBeIgnored: errorCanBeIgnored(e),
				}
				toast.send()
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

type replay struct {
	Type     replayType `json:"type"`
	Offset   string     `json:"offset"`
	Pattern  string     `json:"pattern"`
	SearchID string     `json:"searchId"`
}

type replayType string

const (
	replayFromBeginning replayType = "beginning"
	replayFromTimestamp replayType = "timestamp"
)

func replayTypeStrings() []replayType {
	return []replayType{replayFromBeginning, replayFromTimestamp}
}

// GetReplayTypesList returns the supported replay types
func GetReplayTypesList() []string {
	return []string{
		string(replayFromBeginning),
		string(replayFromTimestamp),
	}
}

// LookupTopic searchs in topic messages, according to the given replay params
func LookupTopic(env base.Environment, topic, params string) (err error) {
	var c *kafka.Consumer

	r := replay{}
	if err = json.Unmarshal([]byte(params), &r); err != nil {
		return
	}

	if r.SearchID == "" {
		err = fmt.Errorf("SearchID parameter must not be empty")
		return
	}

	if !slice.Contains(replayTypeStrings(), r.Type) {
		err = fmt.Errorf("Unsupported Type parameter: %v", r.Type)
		return
	}

	var rx *regexp.Regexp
	if rx, err = regexp.Compile(r.Pattern); err != nil {
		return
	}

	config := cloneConfiguration(env, topic)
	config["enable.auto.commit"] = false
	config["auto.offset.reset"] = "earliest"
	config["go.application.rebalance.enable"] = true
	config["go.events.channel.enable"] = true

	c, err = kafka.NewConsumer(&config)
	if err != nil {
		return
	}

	go func() {
		defer c.Close()

		c.Subscribe(topic, nil)

	replayLoop:
		for {
			select {
			case ev := <-c.Events():
				switch e := ev.(type) {
				case kafka.AssignedPartitions:
					partitions := e.Partitions
					if len(partitions) == 0 {
						log.Info("No partitions assigned")
						continue replayLoop
					}

					switch r.Type {
					case replayFromBeginning:
						log.Info("Replay from beginning, resetting offsets to beginning")

						var s []kafka.TopicPartition
						for _, p := range partitions {
							s = append(s, kafka.TopicPartition{Topic: p.Topic, Partition: p.Partition, Offset: kafka.OffsetBeginning})
						}

						partitions = s
					case replayFromTimestamp:
						log.Infof("Replay from timestamp %s, resetting offsets to that point\n", r.Offset)
						t, err := time.Parse(time.RFC3339Nano, r.Offset)
						if err != nil {
							log.Fatalf("failed to parse replay timestamp %s due to error %v", r.Offset, err)
						}

						var s []kafka.TopicPartition
						for _, p := range partitions {
							s = append(s, kafka.TopicPartition{Topic: p.Topic, Partition: p.Partition, Offset: kafka.Offset(t.UnixNano() / int64(time.Millisecond))})
						}

						if partitions, err = c.OffsetsForTimes(s, 5000); err != nil {
							log.Error(err)
							return
						}
					}

					c.Assign(partitions)
				case kafka.RevokedPartitions:
					c.Unassign()
				case *kafka.Message:
					if rx.MatchString(string(e.Value)) {
						log.Infof("Lookup Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
						if e.Headers != nil {
							log.Infof("Lookup Headers: %v\n", e.Headers)
						}
						sendKafkaMessage(e, r.SearchID)
					}
				case kafka.PartitionEOF:
					log.Info("Reached EOF")
					break replayLoop
				case kafka.Error:
					log.WithFields(log.Fields{"code": e.Code()}).Error(e.String())
					toast := toastMsg{
						Title:        "Lookup Error",
						Message:      e.String(),
						ToastType:    toastError,
						CanBeIgnored: errorCanBeIgnored(e),
					}
					toast.send()
					if e.Code() == kafka.ErrAllBrokersDown ||
						e.Code() == kafka.ErrTimedOut {
						log.Infof("Exiting loop due to error code %v", e.Code())
						break replayLoop
					}
				default:
					log.Warnf("Uncaught event: %v", e)
				}
			}
		}
		c.Unsubscribe()
		log.Infof("Unsubscribed consumer for loop %v", r.SearchID)
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
