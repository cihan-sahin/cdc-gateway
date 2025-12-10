package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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
	brokerFlag := flag.String("broker", "kafka:9092", "Kafka broker address (e.g. localhost:9092 or kafka:9092)")
	topicFlag := flag.String("topic", "dbserver1.inventory.customers", "Kafka topic to produce to")
	countFlag := flag.Int("count", 20, "Number of messages to produce")
	keysFlag := flag.Int("keys", 1, "Number of distinct keys (per-key coalescing test)")
	sleepMsFlag := flag.Int("sleep-ms", 100, "Sleep in ms between messages")
	flipStatusFlag := flag.Bool("flip-status", false, "Toggle status between ACTIVE/INACTIVE to test policy conditions")
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

	log.Printf("Producing test CDC events to Kafka broker=%s topic=%s", broker, topic)
	log.Printf("count=%d, keys=%d, sleep_ms=%d, flip_status=%v",
		*countFlag, *keysFlag, *sleepMsFlag, *flipStatusFlag)

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

	for i := 0; i < *countFlag; i++ {
		id := (i % *keysFlag) + 1
		key := strconv.Itoa(id)

		status := "ACTIVE"
		if *flipStatusFlag && i%5 == 0 {
			status = "INACTIVE"
		}

		env := debeziumEnvelope{}
		env.Payload.Op = "u"
		env.Payload.After = map[string]interface{}{
			"id":     id,
			"name":   fmt.Sprintf("User-%d", id),
			"status": status,
		}
		env.Payload.Source.DB = "dbserver1"
		env.Payload.Source.Schema = "inventory"
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

		log.Printf("Sent CDC event #%d key=%s status=%s cpu=%v",
			i+1, key, status, env.Payload.After["cpu"])

		time.Sleep(time.Duration(*sleepMsFlag) * time.Millisecond)
	}

	log.Println("Done producing test events.")
}
