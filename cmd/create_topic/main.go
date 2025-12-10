package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func main() {
	broker := flag.String("broker", "localhost:9092", "Kafka broker address")
	topic := flag.String("topic", "dbserver1.inventory.customers", "Kafka topic name")
	partitions := flag.Int("partitions", 1, "Number of partitions")
	replicationFactor := flag.Int("replication", 1, "Replication factor")
	flag.Parse()

	if *topic == "" {
		log.Fatal("topic is required (use -topic)")
	}

	log.Printf("Creating topic '%s' on broker %s\n", *topic, *broker)

	conn, err := kafka.Dial("tcp", *broker)
	if err != nil {
		log.Fatalf("failed to dial kafka: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("failed to get controller: %v", err)
	}

	log.Printf("Controller: %v\n", controller)

	controllerAddr := controller.Host + ":" + strconv.Itoa(controller.Port)
	controllerConn, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		log.Fatalf("failed to dial controller (%s): %v", controllerAddr, err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             *topic,
		NumPartitions:     *partitions,
		ReplicationFactor: *replicationFactor,
	})
	if err != nil {
		log.Fatalf("failed to create topic: %v", err)
	}

	log.Printf("Topic '%s' created (partitions=%d, replication=%d)\n",
		*topic, *partitions, *replicationFactor)
}
