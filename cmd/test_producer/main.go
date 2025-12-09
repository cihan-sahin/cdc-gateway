package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	defaultBroker = "localhost:29092"
	defaultTopic  = "cdc.test.customers"
)

type debeziumEnvelope struct {
	Payload struct {
		Op     string                 `json:"op"`
		After  map[string]interface{} `json:"after"`
		Source struct {
			DB     string `json:"db"`
			Schema string `json:"schema"`
			Table  string `json:"table"`
		} `json:"source"`
	} `json:"payload"`
}

func main() {
	broker := defaultBroker
	topic := defaultTopic

	log.Printf("Producing test CDC events to Kafka %s topic %s\n", broker, topic)

	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer w.Close()

	ctx := context.Background()

	key := "1"

	for i := 0; i < 10; i++ {
		env := debeziumEnvelope{}
		env.Payload.Op = "u"
		env.Payload.After = map[string]interface{}{
			"id":     1,
			"name":   fmt.Sprintf("User-%d", i),
			"status": "ACTIVE",
			"cpu":    50 + rand.Intn(50),
		}
		env.Payload.Source.DB = "db"
		env.Payload.Source.Schema = "public"
		env.Payload.Source.Table = "customers"

		valueBytes, err := json.Marshal(env)
		if err != nil {
			log.Fatalf("failed to marshal envelope: %v", err)
		}

		msg := kafka.Message{
			Key:   []byte(key),
			Value: valueBytes,
		}

		if err := w.WriteMessages(ctx, msg); err != nil {
			log.Fatalf("failed to write message: %v", err)
		}

		log.Printf("Sent CDC event #%d for key=%s\n", i+1, key)
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Done producing test events.")
}
