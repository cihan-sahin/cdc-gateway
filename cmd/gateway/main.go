package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfgpkg "github.com/cihan-sahin/cdc-gateway/internal/config"
	"github.com/cihan-sahin/cdc-gateway/internal/kafka"
	"github.com/cihan-sahin/cdc-gateway/internal/metrics"
	"github.com/cihan-sahin/cdc-gateway/internal/pipeline"
	"github.com/cihan-sahin/cdc-gateway/internal/policy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	cfg      *cfgpkg.Config
	policy   policy.Store
	metrics  *metrics.GatewayMetrics
	registry *prometheus.Registry
}

func NewApp(cfg *cfgpkg.Config, store policy.Store, m *metrics.GatewayMetrics, reg *prometheus.Registry) *App {
	return &App{
		cfg:      cfg,
		policy:   store,
		metrics:  m,
		registry: reg,
	}
}

func (a *App) StartHTTP(ctx context.Context) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		in, out := a.metrics.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"events_in": %d, "events_out": %d}`, in, out)
	})

	mux.HandleFunc("/policies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		pols := a.policy.All()
		type respPolicy struct {
			Table         string      `json:"table"`
			Mode          string      `json:"mode"`
			WindowMs      int         `json:"window_ms"`
			MaxBatchSize  int         `json:"max_batch_size"`
			MergeStrategy string      `json:"merge_strategy"`
			TargetTopic   string      `json:"target_topic"`
			Enabled       bool        `json:"enabled"`
			Conditions    interface{} `json:"conditions"`
		}

		out := make([]respPolicy, 0, len(pols))
		for _, p := range pols {
			out = append(out, respPolicy{
				Table:         p.Table,
				Mode:          string(p.Mode),
				WindowMs:      p.WindowMs,
				MaxBatchSize:  p.MaxBatchSize,
				MergeStrategy: string(p.MergeStrategy),
				TargetTopic:   p.TargetTopic,
				Enabled:       p.Enabled,
				Conditions:    p.Conditions,
			})
		}

		b, err := jsonMarshal(out)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"marshal"}`))
			return
		}
		_, _ = w.Write(b)
	})

	mux.Handle("/metrics", promhttp.HandlerFor(a.registry, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:    a.cfg.App.HTTPAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("HTTP server listening on %s", a.cfg.App.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	return srv
}

func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func main() {
	configPath := os.Getenv("CDC_GATEWAY_CONFIG")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := cfgpkg.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Loaded config: kafka brokers=%v, input_topic=%s, group=%s",
		cfg.Kafka.Brokers, cfg.Kafka.InputTopic, cfg.Kafka.ConsumerGroup)

	policyStore, err := cfg.ToPolicyStore()
	if err != nil {
		log.Fatalf("failed to init policy store: %v", err)
	}

	reg := prometheus.NewRegistry()
	m := metrics.NewGatewayMetrics(reg)

	app := NewApp(cfg, policyStore, m, reg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	_ = app.StartHTTP(ctx)

	consumer := kafka.NewKafkaGoConsumer(kafka.KafkaGoConsumerConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.InputTopic,
		GroupID: cfg.Kafka.ConsumerGroup,
	})

	producer := kafka.NewKafkaGoProducer(kafka.KafkaGoProducerConfig{
		Brokers: cfg.Kafka.Brokers,
	})

	router := pipeline.NewRouter(policyStore)
	coalescer := pipeline.NewSimpleCoalescer()

	engine := pipeline.NewEngine(consumer, producer, router, coalescer, m)

	go func() {
		if err := engine.Run(ctx); err != nil {
			log.Printf("engine stopped with error: %v", err)
			stop()
		}
	}()

	log.Printf("cdc-gateway started (env=%s)", cfg.App.Env)

	<-ctx.Done()
	log.Printf("shutting down...")

	_ = consumer.Close()
	_ = producer.Close()
}
