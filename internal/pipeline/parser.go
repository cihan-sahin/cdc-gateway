package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/cihan-sahin/cdc-gateway/internal/kafka"
	"github.com/cihan-sahin/cdc-gateway/pkg/cdcmodel"
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

func ParseDebeziumMessage(msg kafka.Message) (cdcmodel.CDCEvent, error) {
	var env debeziumEnvelope
	if err := json.Unmarshal(msg.Value, &env); err != nil {
		return cdcmodel.CDCEvent{}, fmt.Errorf("unmarshal debezium: %w", err)
	}

	op := cdcmodel.Operation(env.Payload.Op)
	fullTable := fmt.Sprintf("%s.%s.%s", env.Payload.Source.DB, env.Payload.Source.Schema, env.Payload.Source.Table)

	// Key'i Kafka message key'den alıyoruz.
	key := string(msg.Key)

	return cdcmodel.CDCEvent{
		Key:         key,
		Table:       fullTable,
		Op:          op,
		SourceTopic: msg.Topic,
		Payload:     env.Payload.After, // şimdilik "after" ile çalışıyoruz
	}, nil
}
