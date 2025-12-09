package policy

import (
	"fmt"
	"strings"
	"sync"
)

type Store interface {
	GetPolicyForTable(table string) (Policy, bool)
	All() []Policy
	Upsert(p Policy) error
	Delete(table string) error
}

type inMemoryStore struct {
	mu      sync.RWMutex
	byTable map[string]Policy
}

func NewStoreFromRaw(rawPolicies []struct {
	Table         string
	Mode          string
	WindowMs      int
	MaxBatchSize  int
	MergeStrategy string
	TargetTopic   string
	Enabled       bool
	Conditions    []map[string]any
}) (*inMemoryStore, error) {
	store := &inMemoryStore{
		byTable: make(map[string]Policy),
	}

	for _, rp := range rawPolicies {
		pol, err := rawToPolicy(rp)
		if err != nil {
			return nil, err
		}
		store.byTable[pol.Table] = pol
	}

	return store, nil
}

func rawToPolicy(rp struct {
	Table         string
	Mode          string
	WindowMs      int
	MaxBatchSize  int
	MergeStrategy string
	TargetTopic   string
	Enabled       bool
	Conditions    []map[string]any
}) (Policy, error) {
	if strings.TrimSpace(rp.Table) == "" {
		return Policy{}, fmt.Errorf("policy table is required")
	}

	pol := Policy{
		Table:         rp.Table,
		Mode:          Mode(rp.Mode),
		WindowMs:      rp.WindowMs,
		MaxBatchSize:  rp.MaxBatchSize,
		MergeStrategy: MergeStrategy(rp.MergeStrategy),
		TargetTopic:   rp.TargetTopic,
		Enabled:       rp.Enabled,
	}

	for _, rc := range rp.Conditions {
		field, _ := rc["field"].(string)
		opStr, _ := rc["op"].(string)
		val, _ := rc["value"]

		if field == "" || opStr == "" {
			continue
		}

		cond := Condition{
			Field: field,
			Op:    ConditionOp(opStr),
			Value: val,
		}
		pol.Conditions = append(pol.Conditions, cond)
	}

	if err := pol.Validate(); err != nil {
		return Policy{}, fmt.Errorf("invalid policy for table %s: %w", rp.Table, err)
	}

	return pol, nil
}

func (s *inMemoryStore) GetPolicyForTable(table string) (Policy, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.byTable[table]
	return p, ok
}

func (s *inMemoryStore) All() []Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]Policy, 0, len(s.byTable))
	for _, p := range s.byTable {
		res = append(res, p)
	}
	return res
}

func (s *inMemoryStore) Upsert(p Policy) error {
	if err := p.Validate(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byTable[p.Table] = p
	return nil
}

func (s *inMemoryStore) Delete(table string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.byTable, table)
	return nil
}
