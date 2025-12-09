// internal/config/config.go
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Name     string `yaml:"name"`
	Env      string `yaml:"env"`
	HTTPAddr string `yaml:"http_addr"`
}

type KafkaConfig struct {
	Brokers       []string `yaml:"brokers"`
	InputTopic    string   `yaml:"input_topic"`
	ConsumerGroup string   `yaml:"consumer_group"`
}

type RawPolicy struct {
	Table         string                 `yaml:"table"`
	Mode          string                 `yaml:"mode"`
	WindowMs      int                    `yaml:"window_ms"`
	MaxBatchSize  int                    `yaml:"max_batch_size"`
	MergeStrategy string                 `yaml:"merge_strategy"`
	TargetTopic   string                 `yaml:"target_topic"`
	Enabled       bool                   `yaml:"enabled"`
	Conditions    []map[string]any       `yaml:"conditions,omitempty"`
	Extra         map[string]interface{} `yaml:",inline"`
}

type Config struct {
	App      AppConfig   `yaml:"app"`
	Kafka    KafkaConfig `yaml:"kafka"`
	Policies []RawPolicy `yaml:"policies"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	// basit validasyon
	if cfg.App.HTTPAddr == "" {
		cfg.App.HTTPAddr = ":8080"
	}
	if len(cfg.Kafka.Brokers) == 0 {
		return nil, fmt.Errorf("kafka.brokers is required")
	}
	if cfg.Kafka.InputTopic == "" {
		return nil, fmt.Errorf("kafka.input_topic is required")
	}
	if cfg.Kafka.ConsumerGroup == "" {
		return nil, fmt.Errorf("kafka.consumer_group is required")
	}

	return &cfg, nil
}
