package pipeline

import (
	"fmt"
	"sync"
	"time"

	"github.com/cihan-sahin/cdc-gateway/internal/policy"
	"github.com/cihan-sahin/cdc-gateway/pkg/cdcmodel"
)

type bufferedEvent struct {
	event       cdcmodel.CDCEvent
	policy      policy.Policy
	partition   int32
	lastUpdated time.Time
}

type bufferedBatch struct {
	events       []cdcmodel.CDCEvent
	policy       policy.Policy
	partition    int32
	firstEventAt time.Time
	lastUpdated  time.Time
}

type SimpleCoalescer struct {
	mu sync.Mutex

	// last_state buffer'ları: partition -> key -> event
	lastStateBuffers map[int32]map[string]*bufferedEvent

	// micro_batch buffer'ları: partition -> table -> batch
	batchBuffers map[int32]map[string]*bufferedBatch
}

func NewSimpleCoalescer() *SimpleCoalescer {
	return &SimpleCoalescer{
		lastStateBuffers: make(map[int32]map[string]*bufferedEvent),
		batchBuffers:     make(map[int32]map[string]*bufferedBatch),
	}
}

func bufferKey(evt cdcmodel.CDCEvent) string {
	return fmt.Sprintf("%s|%s", evt.Table, evt.Key)
}

func (c *SimpleCoalescer) AddEvent(
	evt cdcmodel.CDCEvent,
	pol policy.Policy,
	partition int32,
	now time.Time,
) ([]CoalescedResult, error) {

	switch pol.Mode {
	case policy.ModePassThrough:
		// Hiç buffer tutma, direkt flush edilmesi için döndür
		return []CoalescedResult{
			{
				Events:       []cdcmodel.CDCEvent{evt},
				Policy:       pol,
				Partition:    partition,
				FlushedAt:    now,
				WindowOpened: now,
			},
		}, nil

	case policy.ModeLastState:
		return c.addLastState(evt, pol, partition, now)

	case policy.ModeMicroBatch:
		return c.addMicroBatch(evt, pol, partition, now)

	default:
		// Tanınmayan mode → passthrough
		return []CoalescedResult{
			{
				Events:       []cdcmodel.CDCEvent{evt},
				Policy:       pol,
				Partition:    partition,
				FlushedAt:    now,
				WindowOpened: now,
			},
		}, nil
	}
}

func (c *SimpleCoalescer) addLastState(
	evt cdcmodel.CDCEvent,
	pol policy.Policy,
	partition int32,
	now time.Time,
) ([]CoalescedResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.lastStateBuffers[partition]; !ok {
		c.lastStateBuffers[partition] = make(map[string]*bufferedEvent)
	}

	k := bufferKey(evt)
	if existing, ok := c.lastStateBuffers[partition][k]; ok {
		// merge_strategy'ye göre yeni event'i oluştur
		merged := MergeEvent(existing.event, evt, pol.MergeStrategy)
		existing.event = merged
		existing.lastUpdated = now
	} else {
		c.lastStateBuffers[partition][k] = &bufferedEvent{
			event:       evt,
			policy:      pol,
			partition:   partition,
			lastUpdated: now,
		}
	}

	// last_state için flush zaman bazlı yapılacak, burada dönmüyoruz
	return nil, nil
}

func (c *SimpleCoalescer) addMicroBatch(
	evt cdcmodel.CDCEvent,
	pol policy.Policy,
	partition int32,
	now time.Time,
) ([]CoalescedResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.batchBuffers[partition]; !ok {
		c.batchBuffers[partition] = make(map[string]*bufferedBatch)
	}

	tableKey := pol.Table
	batch, ok := c.batchBuffers[partition][tableKey]
	if !ok {
		batch = &bufferedBatch{
			policy:       pol,
			partition:    partition,
			firstEventAt: now,
			lastUpdated:  now,
			events:       make([]cdcmodel.CDCEvent, 0, max(1, pol.MaxBatchSize)),
		}
		c.batchBuffers[partition][tableKey] = batch
	}

	batch.events = append(batch.events, evt)
	batch.lastUpdated = now

	// Eğer max_batch_size konmuşsa ve doldurduysak hemen flush edelim
	if pol.MaxBatchSize > 0 && len(batch.events) >= pol.MaxBatchSize {
		res := CoalescedResult{
			Events:       batch.events,
			Policy:       pol,
			Partition:    partition,
			FlushedAt:    now,
			WindowOpened: batch.firstEventAt,
		}
		// batch'i sıfırla
		batch.events = nil
		batch.firstEventAt = now
		return []CoalescedResult{res}, nil
	}

	// Aksi halde FlushDue zamanında kontrol edilecek
	return nil, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (c *SimpleCoalescer) FlushDue(now time.Time) ([]CoalescedResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var results []CoalescedResult

	// LAST_STATE için: window_ms'e göre flush
	for partition, buf := range c.lastStateBuffers {
		for k, be := range buf {
			win := time.Duration(be.policy.WindowMs) * time.Millisecond
			if win <= 0 {
				continue
			}
			if now.Sub(be.lastUpdated) >= win {
				results = append(results, CoalescedResult{
					Events:       []cdcmodel.CDCEvent{be.event},
					Policy:       be.policy,
					Partition:    partition,
					FlushedAt:    now,
					WindowOpened: be.lastUpdated.Add(-win),
				})
				delete(buf, k)
			}
		}
	}

	// MICRO_BATCH için: window_ms'e göre flush
	for partition, tableBuf := range c.batchBuffers {
		for table, batch := range tableBuf {
			if len(batch.events) == 0 {
				continue
			}
			win := time.Duration(batch.policy.WindowMs) * time.Millisecond
			if win <= 0 {
				continue
			}
			if now.Sub(batch.firstEventAt) >= win {
				results = append(results, CoalescedResult{
					Events:       batch.events,
					Policy:       batch.policy,
					Partition:    partition,
					FlushedAt:    now,
					WindowOpened: batch.firstEventAt,
				})
				// batch temizle
				batch.events = nil
				batch.firstEventAt = now
				// istersen map'ten tamamen silebilirsin:
				// delete(tableBuf, table)
				_ = table // sadece örnek
			}
		}
	}

	return results, nil
}
