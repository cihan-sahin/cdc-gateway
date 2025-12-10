package metrics

import (
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type GatewayMetrics struct {
	EventsIn  atomic.Int64
	EventsOut atomic.Int64

	eventsInCounter  prometheus.Counter
	eventsOutCounter prometheus.Counter
	batchSizeHist    prometheus.Histogram
	flushLatencyHist prometheus.Histogram
}

func NewGatewayMetrics(reg prometheus.Registerer) *GatewayMetrics {
	mEventsIn := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cdc_gateway_events_in_total",
		Help: "Total number of CDC events received from Kafka",
	})
	mEventsOut := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cdc_gateway_events_out_total",
		Help: "Total number of events emitted after coalescing/batching",
	})
	batchSizeHist := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "cdc_gateway_batch_size",
		Help:    "Number of events in each flushed batch",
		Buckets: prometheus.ExponentialBuckets(1, 2, 8),
	})
	flushLatencyHist := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "cdc_gateway_flush_latency_seconds",
		Help:    "Latency between first event in window and flush time",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
	})

	reg.MustRegister(mEventsIn, mEventsOut, batchSizeHist, flushLatencyHist)

	return &GatewayMetrics{
		eventsInCounter:  mEventsIn,
		eventsOutCounter: mEventsOut,
		batchSizeHist:    batchSizeHist,
		flushLatencyHist: flushLatencyHist,
	}
}

func (m *GatewayMetrics) IncEventsIn() {
	m.EventsIn.Add(1)
	if m.eventsInCounter != nil {
		m.eventsInCounter.Inc()
	}
}

func (m *GatewayMetrics) IncEventsOut(n int64) {
	m.EventsOut.Add(n)
	if m.eventsOutCounter != nil {
		m.eventsOutCounter.Add(float64(n))
	}
}

func (m *GatewayMetrics) ObserveBatch(size int, windowOpened, flushedAt time.Time) {
	if m.batchSizeHist != nil && size > 0 {
		m.batchSizeHist.Observe(float64(size))
	}
	if m.flushLatencyHist != nil && !windowOpened.IsZero() {
		lat := flushedAt.Sub(windowOpened).Seconds()
		if lat >= 0 {
			m.flushLatencyHist.Observe(lat)
		}
	}
}

func (m *GatewayMetrics) Snapshot() (eventsIn, eventsOut int64) {
	return m.EventsIn.Load(), m.EventsOut.Load()
}
