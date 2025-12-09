package pipeline

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/cihan-sahin/cdc-gateway/internal/kafka"
	"github.com/cihan-sahin/cdc-gateway/internal/metrics"
	"github.com/cihan-sahin/cdc-gateway/internal/policy"
	"github.com/cihan-sahin/cdc-gateway/pkg/cdcmodel"
)

type Engine struct {
	consumer  kafka.Consumer
	producer  kafka.Producer
	router    *Router
	coalescer Coalescer
	metrics   *metrics.GatewayMetrics
}

func NewEngine(consumer kafka.Consumer, producer kafka.Producer, router *Router, c Coalescer, m *metrics.GatewayMetrics) *Engine {
	return &Engine{
		consumer:  consumer,
		producer:  producer,
		router:    router,
		coalescer: c,
		metrics:   m,
	}
}

func (e *Engine) Run(ctx context.Context) error {
	// Flush ticker: periyodik olarak buffer'ı kontrol et
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				results, err := e.coalescer.FlushDue(t)
				if err != nil {
					log.Printf("FlushDue error: %v", err)
					continue
				}
				if err := e.handleResults(ctx, results); err != nil {
					log.Printf("handleResults error: %v", err)
				}
			}
		}
	}()

	// Consumer loop
	return e.consumer.Consume(ctx, e.handleMessage)
}

func (e *Engine) handleMessage(ctx context.Context, msg kafka.Message) error {
	e.metrics.IncEventsIn()

	evt, err := ParseDebeziumMessage(msg)
	if err != nil {
		log.Printf("parse debezium error: %v", err)
		return nil // drop, ama istersen DLQ yaparız
	}

	routed, ok := e.router.Route(evt)
	if !ok {
		// Policy yok veya conditions uyuşmadı → raw event'i passthrough gibi davranalım
		return e.sendRaw(ctx, evt)
	}

	now := time.Now()
	results, err := e.coalescer.AddEvent(routed.Event, routed.Policy, int32(msg.Partition), now)
	if err != nil {
		log.Printf("AddEvent error: %v", err)
		return nil
	}

	if len(results) > 0 {
		if err := e.handleResults(ctx, results); err != nil {
			log.Printf("handleResults error: %v", err)
		}
	}

	return nil
}

func (e *Engine) handleResults(ctx context.Context, results []CoalescedResult) error {
	for _, res := range results {
		// Prometheus’a batch metriklerini yaz
		e.metrics.ObserveBatch(len(res.Events), res.WindowOpened, res.FlushedAt)

		switch res.Policy.Mode {
		case policy.ModeMicroBatch:
			if err := e.sendBatch(ctx, res); err != nil {
				log.Printf("sendBatch error: %v", err)
			}
			e.metrics.IncEventsOut(int64(len(res.Events)))
		default:
			for _, evt := range res.Events {
				if err := e.sendCoalesced(ctx, evt, res.Policy, res.Partition); err != nil {
					log.Printf("sendCoalesced error: %v", err)
				}
				e.metrics.IncEventsOut(1)
			}
		}
	}
	return nil
}

func (e *Engine) sendRaw(ctx context.Context, evt cdcmodel.CDCEvent) error {
	// Raw event'i, istersen aynı topic'e veya ayrı "raw" topic'e basabilirsin.
	// Şimdilik: orijinal topic'e JSON serialize ederek atalım.
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	return e.producer.Send(ctx, kafka.Message{
		Topic: evt.SourceTopic,
		Key:   []byte(evt.Key),
		Value: b,
	})
}

func (e *Engine) sendCoalesced(ctx context.Context, evt cdcmodel.CDCEvent, pol policy.Policy, partition int32) error {
	b, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	topic := pol.TargetTopic
	if topic == "" {
		// yedek: orijinal topic
		topic = evt.SourceTopic
	}
	return e.producer.Send(ctx, kafka.Message{
		Topic:     topic,
		Partition: int(partition),
		Key:       []byte(evt.Key),
		Value:     b,
	})
}

func (e *Engine) sendBatch(ctx context.Context, res CoalescedResult) error {
	// batch payload formatını basit tutalım:
	// { "table": "...", "events": [ {...}, {...} ] }
	payload := map[string]interface{}{
		"table":  res.Policy.Table,
		"events": res.Events,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	topic := res.Policy.TargetTopic
	if topic == "" && len(res.Events) > 0 {
		topic = res.Events[0].SourceTopic
	}

	return e.producer.Send(ctx, kafka.Message{
		Topic:     topic,
		Partition: int(res.Partition),
		Key:       nil, // batch key'i yok, consumers kendisi handle edecek
		Value:     b,
	})
}
