package tests

import (
	"testing"
	"time"

	"github.com/cihan-sahin/cdc-gateway/internal/pipeline"
	"github.com/cihan-sahin/cdc-gateway/internal/policy"
	"github.com/cihan-sahin/cdc-gateway/pkg/cdcmodel"
)

func TestSimpleCoalescer_LastStateFlush(t *testing.T) {
	c := pipeline.NewSimpleCoalescer()

	pol := policy.Policy{
		Table:         "db.public.users",
		Mode:          policy.ModeLastState,
		WindowMs:      100, // 100ms window
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

	// AddEvent sadece buffer'a koyacak
	results, err := c.AddEvent(evt1, pol, 0, now)
	if err != nil {
		t.Fatalf("AddEvent error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no immediate flush for last_state, got %d", len(results))
	}

	// 200ms sonra FlushDue çağır
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
