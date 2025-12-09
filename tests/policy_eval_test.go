package tests

import (
	"testing"

	"github.com/cihan-sahin/cdc-gateway/internal/policy"
)

func TestPolicyMatchesPayload_WithConditions(t *testing.T) {
	p := policy.Policy{
		Table: "public.device_metrics",
		Mode:  policy.ModeLastState,
		Conditions: []policy.Condition{
			{
				Field: "type",
				Op:    policy.OpEq,
				Value: "STATUS_UPDATE",
			},
			{
				Field: "cpu",
				Op:    policy.OpGt,
				Value: 80,
			},
		},
	}

	payload := map[string]interface{}{
		"type": "STATUS_UPDATE",
		"cpu":  92,
	}

	if !p.MatchesPayload(payload) {
		t.Fatalf("expected payload to match policy conditions")
	}

	payload2 := map[string]interface{}{
		"type": "STATUS_UPDATE",
		"cpu":  50,
	}

	if p.MatchesPayload(payload2) {
		t.Fatalf("expected payload NOT to match policy (cpu too low)")
	}
}
