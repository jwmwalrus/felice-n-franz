package kafkian

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"time"

	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/felice-n-franz/internal/base"
	rtc "github.com/jwmwalrus/rtcycler"
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

// GetReplayTypesList returns the supported replay types.
func GetReplayTypesList() []string {
	return []string{
		string(replayFromBeginning),
		string(replayFromTimestamp),
	}
}

// LookupTopic searchs in topic messages, according to the given replay params.
func LookupTopic(env base.Environment, topic, params string) (err error) {
	var c *kafka.Consumer

	r := replay{}
	if err = json.Unmarshal([]byte(params), &r); err != nil {
		return
	}

	if r.SearchID == "" {
		err = fmt.Errorf("searchID parameter must not be empty")
		return
	}

	if !slices.Contains(replayTypeStrings(), r.Type) {
		err = fmt.Errorf("unsupported Type parameter: %v", r.Type)
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
				slog.Info("Consumer is not registered anymore!")
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
						slog.Info("No partitions assigned")
						continue replayLoop
					}

					switch r.Type {
					case replayFromBeginning:
						slog.Info("Replay from beginning, resetting offsets to beginning")

						var s []kafka.TopicPartition
						for _, p := range partitions {
							s = append(s, kafka.TopicPartition{Topic: p.Topic, Partition: p.Partition, Offset: kafka.OffsetBeginning})
						}

						partitions = s
					case replayFromTimestamp:
						slog.Info("Replay from timestamp, resetting offsets to that point", "offset", r.Offset)
						t, err := time.Parse(time.RFC3339Nano, r.Offset)
						if err != nil {
							rtc.With(
								"offset", r.Offset,
								"error", err,
							).Fatal("Failed to parse replay timestamp")
						}

						var s []kafka.TopicPartition
						for _, p := range partitions {
							s = append(s, kafka.TopicPartition{Topic: p.Topic, Partition: p.Partition, Offset: kafka.Offset(t.UnixNano() / int64(time.Millisecond))})
						}

						if partitions, err = c.OffsetsForTimes(s, 5000); err != nil {
							slog.Error("Failed to get partitions", "error", err)
							return
						}
					}

					c.Assign(partitions)
				case kafka.RevokedPartitions:
					c.Unassign()
				case *kafka.Message:
					if rx.MatchString(string(e.Value)) {
						slog.With(
							"topic-partition", e.TopicPartition,
							"value", string(e.Value),
						).Info("Lookup Message")
						if e.Headers != nil {
							slog.Info("Lookup Headers", "headers", e.Headers)
						}
						sendKafkaMessage(e, r.SearchID)
					}
				case kafka.PartitionEOF:
					slog.Info("Reached EOF")
					break replayLoop
				case kafka.Error:
					slog.With(
						"code", e.Code(),
						"error", e,
					).Error("Kafka error")
					toast := toastMsg{
						ToastType:    toastError,
						Title:        "Lookup Error",
						Message:      e.String(),
						CanBeIgnored: errorCanBeIgnored(e),
					}
					toast.send()
					if e.Code() == kafka.ErrAllBrokersDown ||
						e.Code() == kafka.ErrTimedOut {
						slog.Info("Exiting loop due to timeout", "code", e.Code())
						break replayLoop
					}
				default:
					slog.Warn("Ignoring Kafka event", "event", e)
				}
			}
		}
		c.Unsubscribe()
		slog.Info("Unsubscribed consumer for loop", "searchID", r.SearchID)
	}()

	return
}
