package tests

import (
	"testing"
	"time"

	"github.com/cihan-sahin/cdc-gateway/internal/policy"
	"github.com/cihan-sahin/cdc-gateway/pkg/cdcmodel"
	"github.com/cihan-sahin/cdc-gateway/internal/pipeline"
)

func TestSimpleCoalescer_LastStateFlushByWindow(t *testing.T) {
	c := pipeline.NewSimpleCoalescer()

	pol := policy.Policy{
		Table:         "db.public.users",
		Mode:          policy.ModeLastState,
		WindowMs:      100,
		MaxBatchSize:  0,
		MergeStrategy: policy.MergeReplace,
		TargetTopic:   "users.optimized",
		Enabled:       true,
	}

	now := time.Now()

	evt1 := cdcmodel.CDCEvent{
		Key:   "1",
		Table: "db.public.users",
		Op:    cdcmodel.OpUpdate,
		Payload: map[string]interface{}{
			"name": "Alice",
		},
	}

	results, err := c.AddEvent(evt1, pol, 0, now)
	if err != nil {
		t.Fatalf("AddEvent error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no immediate flush for last_state, got %d", len(results))
	}

	later := now.Add(200 * time.Millisecond)
	results, err = c.FlushDue(later)
	if err != nil {
		t.Fatalf("FlushDue error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 flush result, got %d", len(results))
	}
	if len(results[0].Events) != 1 {
		t.Fatalf("expected 1 event in result, got %d", len(results[0].Events))
	}
	if results[0].Events[0].Payload["name"] != "Alice" {
		t.Fatalf("unexpected payload: %+v", results[0].Events[0].Payload)
	}
}

func TestSimpleCoalescer_MicroBatch_FlushByBatchSize(t *testing.T) {
	c := pipeline.NewSimpleCoalescer()

	pol := policy.Policy{
		Table:         "db.public.device_metrics",
		Mode:          policy.ModeMicroBatch,
		WindowMs:      1000,
		MaxBatchSize:  3,
		MergeStrategy: policy.MergeReplace,
		TargetTopic:   "device_metrics.batched",
		Enabled:       true,
	}

	now := time.Now()

	for i := 0; i < 3; i++ {
		evt := cdcmodel.CDCEvent{
			Key:   "",
			Table: "db.public.device_metrics",
			Op:    cdcmodel.OpUpdate,
			Payload: map[string]interface{}{
				"seq": i,
			},
		}

		results, err := c.AddEvent(evt, pol, 0, now.Add(time.Duration(i)*10*time.Millisecond))
		if err != nil {
			t.Fatalf("AddEvent error: %v", err)
		}

		if i < 2 {
			if len(results) != 0 {
				t.Fatalf("expected no flush at i=%d, got %d results", i, len(results))
			}
		} else {
			if len(results) != 1 {
				t.Fatalf("expected 1 flush result at i=%d, got %d", i, len(results))
			}
			if len(results[0].Events) != 3 {
				t.Fatalf("expected 3 events in batch, got %d", len(results[0].Events))
			}
		}
	}
}

func TestSimpleCoalescer_MicroBatch_FlushByWindow(t *testing.T) {
	c := pipeline.NewSimpleCoalescer()

	pol := policy.Policy{
		Table:         "db.public.logs",
		Mode:          policy.ModeMicroBatch,
		WindowMs:      100, // 100ms window
		MaxBatchSize:  0,   // batch size limiti yok, sadece window'a bak
		MergeStrategy: policy.MergeReplace,
		TargetTopic:   "logs.batched",
		Enabled:       true,
	}

	now := time.Now()

	evt1 := cdcmodel.CDCEvent{
		Table: "db.public.logs",
		Op:    cdcmodel.OpCreate,
		Payload: map[string]interface{}{
			"msg": "first",
		},
	}
	evt2 := cdcmodel.CDCEvent{
		Table: "db.public.logs",
		Op:    cdcmodel.OpCreate,
		Payload: map[string]interface{}{
			"msg": "second",
		},
	}

	if res, err := c.AddEvent(evt1, pol, 0, now); err != nil {
		t.Fatalf("AddEvent error: %v", err)
	} else if len(res) != 0 {
		t.Fatalf("expected no flush for first event, got %d", len(res))
	}

	if res, err := c.AddEvent(evt2, pol, 0, now.Add(50*time.Millisecond)); err != nil {
		t.Fatalf("AddEvent error: %v", err)
	} else if len(res) != 0 {
		t.Fatalf("expected no flush for second event, got %d", len(res))
	}

	later := now.Add(200 * time.Millisecond)
	results, err := c.FlushDue(later)
	if err != nil {
		t.Fatalf("FlushDue error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 flush result, got %d", len(results))
	}
	if len(results[0].Events) != 2 {
		t.Fatalf("expected 2 events in batch, got %d", len(results[0].Events))
	}
}

func TestMergeEvent_Replace(t *testing.T) {
	prev := cdcmodel.CDCEvent{
		Payload: map[string]interface{}{
			"name": "Alice",
			"age":  25,
		},
	}
	next := cdcmodel.CDCEvent{
		Payload: map[string]interface{}{
			"name": "Alice2",
		},
	}

	merged := pipeline.MergeEvent(prev, next, policy.MergeReplace)
	if merged.Payload["name"] != "Alice2" {
		t.Fatalf("expected name=Alice2, got=%v", merged.Payload["name"])
	}
	if _, ok := merged.Payload["age"]; ok {
		t.Fatalf("expected age field to be gone in replace strategy")
	}
}

func TestMergeEvent_MergeFields(t *testing.T) {
	prev := cdcmodel.CDCEvent{
		Payload: map[string]interface{}{
			"name": "Alice",
			"age":  25,
		},
	}
	next := cdcmodel.CDCEvent{
		Payload: map[string]interface{}{
			"name": "Alice2",
		},
	}

	merged := pipeline.MergeEvent(prev, next, policy.MergeMergeFields)
	if merged.Payload["name"] != "Alice2" {
		t.Fatalf("expected name=Alice2, got=%v", merged.Payload["name"])
	}
	if merged.Payload["age"] != 25.0 && merged.Payload["age"] != 25 {
		t.Fatalf("expected age=25, got=%v", merged.Payload["age"])
	}
}

func TestMergeEvent_AppendHistory(t *testing.T) {
	prev := cdcmodel.CDCEvent{
		Payload: map[string]interface{}{
			"status": "pending",
		},
	}
	next := cdcmodel.CDCEvent{
		Payload: map[string]interface{}{
			"status": "completed",
		},
	}

	merged := pipeline.MergeEvent(prev, next, policy.MergeAppendHistory)
	h, ok := merged.Payload["_history"]
	if !ok {
		t.Fatalf("expected _history field in payload")
	}

	hSlice, ok := h.([]map[string]interface{})
	if !ok {
		t.Fatalf("expected _history to be []map[string]interface{}, got=%T", h)
	}
	if len(hSlice) != 1 {
		t.Fatalf("expected history length=1, got=%d", len(hSlice))
	}

	entry := hSlice[0]
	data, ok := entry["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected history entry data to be map[string]interface{}, got=%T", entry["data"])
	}
	if data["status"] != "pending" {
		t.Fatalf("expected previous status in history to be 'pending', got=%v", data["status"])
	}
	if merged.Payload["status"] != "completed" {
		t.Fatalf("expected current status to be 'completed', got=%v", merged.Payload["status"])
	}
}
