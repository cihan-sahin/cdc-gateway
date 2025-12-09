package pipeline

import (
	"time"

	"github.com/cihan-sahin/cdc-gateway/internal/policy"
	"github.com/cihan-sahin/cdc-gateway/pkg/cdcmodel"
)

// MergeEvent: aynı key için ardışık eventleri, policy'ye göre birleştirir.
func MergeEvent(prev, next cdcmodel.CDCEvent, strat policy.MergeStrategy) cdcmodel.CDCEvent {
	switch strat {
	case policy.MergeReplace:
		// En basit: son gelen kazanır
		return next

	case policy.MergeMergeFields:
		return mergeFields(prev, next)

	case policy.MergeAppendHistory:
		return appendHistory(prev, next)

	default:
		return next
	}
}

// merge_fields: önceki payload'taki alanları koru, yeni gelen alanları override et.
func mergeFields(prev, next cdcmodel.CDCEvent) cdcmodel.CDCEvent {
	merged := make(map[string]interface{}, len(prev.Payload)+len(next.Payload))
	for k, v := range prev.Payload {
		merged[k] = v
	}
	for k, v := range next.Payload {
		merged[k] = v
	}
	next.Payload = merged
	return next
}

// append_history:
// payload:
//
//	current: {...son state...}
//	history: [
//	  { at: <time>, data: {...} },
//	  ...
//	]
func appendHistory(prev, next cdcmodel.CDCEvent) cdcmodel.CDCEvent {
	historyAny, ok := prev.Payload["_history"]
	var history []map[string]interface{}
	if ok {
		if slice, ok := historyAny.([]map[string]interface{}); ok {
			history = slice
		}
	}

	// Bir önceki state'i history'ye at
	entry := map[string]interface{}{
		"at":   time.Now().UTC().Format(time.RFC3339Nano),
		"data": prev.Payload,
	}
	history = append(history, entry)

	// Yeni payload'a history'yi ekle
	if next.Payload == nil {
		next.Payload = make(map[string]interface{})
	}
	next.Payload["_history"] = history

	return next
}
