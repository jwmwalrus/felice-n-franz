package kafkian

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/jwmwalrus/felice-n-franz/internal/base"
	"github.com/jwmwalrus/onerror"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

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

	if !slices.Contains(replayTypeStrings(), r.Type) {
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

	reg := register(c, env.Name, topic, r.SearchID)

	go func() {
		defer c.Close()

		c.Subscribe(topic, nil)

		toast := toastMsg{
			ToastType: toastInfo,
			Title:     "From Lookup",
			Message:   "Replay started for searchId " + r.SearchID,
		}
		toast.send()

	replayLoop:
		for {
			select {
			case <-reg.quit:
				log.Info("Consumer is not registered anymore!")
				err := c.Unsubscribe()
				onerror.Warn(err)
				toast := toastMsg{
					ToastType: toastInfo,
					Title:     "From Lookup",
					Message:   "Replay stopped for searchId " + r.SearchID,
				}
				toast.send()
				break replayLoop
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
						ToastType:    toastError,
						Title:        "Lookup Error",
						Message:      e.String(),
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
