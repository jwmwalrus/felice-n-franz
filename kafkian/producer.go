package kafkian

import (
	"github.com/jwmwalrus/felice-n-franz/base"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func ProduceMessage(env base.Environment, topic string, payload []byte) (err error) {

	var p *kafka.Producer
	p, err = kafka.NewProducer(&env.Configuration)
	if err != nil {
		return
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				toast := toastMsg{
					ToastType: "info",
					Title:     "From Producer",
					Message:   "Message delivered",
					Topic:     *ev.TopicPartition.Topic,
					Partition: ev.TopicPartition.Partition,
					Offset:    ev.TopicPartition.Offset,
				}
				if ev.TopicPartition.Error != nil {
					toast.ToastType = "error"
					toast.Message = "Delivery failed"
				}
				sendToast(toast)
			}
		}
	}()

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, nil)

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
	return
}
