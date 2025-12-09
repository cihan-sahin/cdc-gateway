package kafka

import "time"

type Message struct {
	Topic     string
	Partition int
	Offset    int64
	Key       []byte
	Value     []byte
	Timestamp time.Time
}
