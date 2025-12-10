package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	defaultBroker = "localhost:9092"
	defaultTopic  = "dbserver1.inventory.customers"
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
	brokerFlag := flag.String("broker", "", "Kafka broker address (e.g. localhost:9092 or kafka:9092)")
	topicFlag := flag.String("topic", "", "Kafka topic to produce to")
	countFlag := flag.Int("count", 10, "Number of messages to produce")
	flag.Parse()

	broker := *brokerFlag
	if broker == "" {
		if env := os.Getenv("KAFKA_BROKER"); env != "" {
			broker = env
		} else {
			broker = defaultBroker
		}
	}

	topic := *topicFlag
	if topic == "" {
		if env := os.Getenv("KAFKA_TOPIC"); env != "" {
			topic = env
		} else {
			topic = defaultTopic
		}
	}

	log.Printf("Producing test CDC events to Kafka broker=%s topic=%s\n", broker, topic)

	w := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer func() {
		if err := w.Close(); err != nil {
			log.Printf("failed to close writer: %v", err)
		}
	}()

	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())
	key := "1"

	for i := 0; i < *countFlag; i++ {
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
