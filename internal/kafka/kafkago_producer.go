// internal/kafka/kafkago_producer.go
package kafka

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

type KafkaGoProducerConfig struct {
	Brokers []string
	// İstersen default topic vs. ekleyebilirsin
}

type KafkaGoProducer struct {
	w *kafkago.Writer
}

func NewKafkaGoProducer(cfg KafkaGoProducerConfig) *KafkaGoProducer {
	w := &kafkago.Writer{
		Addr:         kafkago.TCP(cfg.Brokers...),
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequireOne,
		Async:        false,
	}
	return &KafkaGoProducer{w: w}
}

func (p *KafkaGoProducer) Send(ctx context.Context, msg Message) error {
	return p.w.WriteMessages(ctx, kafkago.Message{
		Topic:     msg.Topic,
		Partition: msg.Partition, // genelde 0 bırakılır, balancer seçer
		Key:       msg.Key,
		Value:     msg.Value,
		Time:      msg.Timestamp,
	})
}

func (p *KafkaGoProducer) Close() error {
	return p.w.Close()
}
