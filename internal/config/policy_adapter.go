package config

import "github.com/cihan-sahin/cdc-gateway/internal/policy"

func (c *Config) ToPolicyStore() (policy.Store, error) {
	raw := make([]struct {
		Table         string
		Mode          string
		WindowMs      int
		MaxBatchSize  int
		MergeStrategy string
		TargetTopic   string
		Enabled       bool
		Conditions    []map[string]any
	}, 0, len(c.Policies))

	for _, p := range c.Policies {
		raw = append(raw, struct {
			Table         string
			Mode          string
			WindowMs      int
			MaxBatchSize  int
			MergeStrategy string
			TargetTopic   string
			Enabled       bool
			Conditions    []map[string]any
		}{
			Table:         p.Table,
			Mode:          p.Mode,
			WindowMs:      p.WindowMs,
			MaxBatchSize:  p.MaxBatchSize,
			MergeStrategy: p.MergeStrategy,
			TargetTopic:   p.TargetTopic,
			Enabled:       p.Enabled,
			Conditions:    p.Conditions,
		})
	}

	return policy.NewStoreFromRaw(raw)
}
