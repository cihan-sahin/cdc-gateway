package policy

import (
	"fmt"
)

func (c Condition) Matches(value interface{}) bool {
	switch c.Op {
	case OpExists:
		return value != nil
	case OpMissing:
		return value == nil
	}

	if value == nil {
		return false
	}

	switch v := value.(type) {
	case string:
		return c.matchString(v)
	case float64, float32, int, int64, int32:
		num, ok := toFloat64(v)
		if !ok {
			return false
		}
		return c.matchNumber(num)
	case bool:
		return c.matchBool(v)
	default:
		return false
	}
}

func (c Condition) matchString(v string) bool {
	target, ok := c.Value.(string)
	if !ok {
		return false
	}
	switch c.Op {
	case OpEq:
		return v == target
	case OpNe:
		return v != target
	case OpIn, OpNotIn:
		list, ok := toStringSlice(c.Value)
		if !ok {
			return false
		}
		contains := stringInSlice(v, list)
		if c.Op == OpIn {
			return contains
		}
		return !contains
	default:
		return false
	}
}

func (c Condition) matchNumber(v float64) bool {
	target, ok := toFloat64(c.Value)
	if !ok {
		return false
	}
	switch c.Op {
	case OpEq:
		return v == target
	case OpNe:
		return v != target
	case OpGt:
		return v > target
	case OpGte:
		return v >= target
	case OpLt:
		return v < target
	case OpLte:
		return v <= target
	default:
		return false
	}
}

func (c Condition) matchBool(v bool) bool {
	target, ok := c.Value.(bool)
	if !ok {
		return false
	}
	switch c.Op {
	case OpEq:
		return v == target
	case OpNe:
		return v != target
	default:
		return false
	}
}

func toFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case int32:
		return float64(n), true
	default:
		return 0, false
	}
}

func toStringSlice(v interface{}) ([]string, bool) {
	switch s := v.(type) {
	case []string:
		return s, true
	case []interface{}:
		res := make([]string, 0, len(s))
		for _, item := range s {
			str, ok := item.(string)
			if !ok {
				return nil, false
			}
			res = append(res, str)
		}
		return res, true
	default:
		return nil, false
	}
}

func stringInSlice(v string, list []string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func (p Policy) MatchesPayload(payload map[string]interface{}) bool {
	if len(p.Conditions) == 0 {
		return true
	}
	for _, cond := range p.Conditions {
		val, _ := payload[cond.Field]
		if !cond.Matches(val) {
			return false
		}
	}
	return true
}

func (p Policy) Validate() error {
	if p.Table == "" {
		return fmt.Errorf("table is required")
	}
	if p.Mode == "" {
		return fmt.Errorf("mode is required")
	}
	return nil
}
