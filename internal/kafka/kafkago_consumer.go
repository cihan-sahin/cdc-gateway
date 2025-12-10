package kafka

import (
	"context"
	"errors"
	"log"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type KafkaGoConsumerConfig struct {
	Brokers               []string
	Topic                 string
	GroupID               string
	MinBytes              int
	MaxBytes              int
	CommitIntervalSeconds int
}

type KafkaGoConsumer struct {
	reader *kafkago.Reader
}

func NewKafkaGoConsumer(cfg KafkaGoConsumerConfig) *KafkaGoConsumer {
	if cfg.MinBytes == 0 {
		cfg.MinBytes = 1e3
	}
	if cfg.MaxBytes == 0 {
		cfg.MaxBytes = 10e6
	}
	if cfg.CommitIntervalSeconds == 0 {
		cfg.CommitIntervalSeconds = 1
	}

	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topic,
		GroupID:        cfg.GroupID,
		MinBytes:       cfg.MinBytes,
		MaxBytes:       cfg.MaxBytes,
		CommitInterval: (time.Duration(cfg.CommitIntervalSeconds) * time.Second),
	})

	return &KafkaGoConsumer{reader: r}
}

func (c *KafkaGoConsumer) Consume(ctx context.Context, handler func(ctx context.Context, msg Message) error) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		msg := Message{
			Topic:     m.Topic,
			Partition: m.Partition,
			Offset:    m.Offset,
			Key:       m.Key,
			Value:     m.Value,
			Timestamp: m.Time,
		}

		if err := handler(ctx, msg); err != nil {
			log.Printf("handler error (partition=%d offset=%d): %v", m.Partition, m.Offset, err)
		}
	}
}

func (c *KafkaGoConsumer) Close() error {
	return c.reader.Close()
}
