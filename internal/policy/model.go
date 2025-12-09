package policy

type Mode string

const (
	ModePassThrough Mode = "passthrough"
	ModeLastState   Mode = "last_state"
	ModeMicroBatch  Mode = "micro_batch"
)

type MergeStrategy string

const (
	MergeReplace       MergeStrategy = "replace"
	MergeMergeFields   MergeStrategy = "merge_fields"
	MergeAppendHistory MergeStrategy = "append_history"
)

type ConditionOp string

const (
	OpEq      ConditionOp = "eq"
	OpNe      ConditionOp = "ne"
	OpGt      ConditionOp = "gt"
	OpGte     ConditionOp = "gte"
	OpLt      ConditionOp = "lt"
	OpLte     ConditionOp = "lte"
	OpIn      ConditionOp = "in"
	OpNotIn   ConditionOp = "not_in"
	OpExists  ConditionOp = "exists"
	OpMissing ConditionOp = "missing"
)

type Condition struct {
	Field string      `json:"field" yaml:"field"`
	Op    ConditionOp `json:"op" yaml:"op"`
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
}

type Policy struct {
	Table         string        `json:"table" yaml:"table"`
	Mode          Mode          `json:"mode" yaml:"mode"`
	WindowMs      int           `json:"window_ms" yaml:"window_ms"`
	MaxBatchSize  int           `json:"max_batch_size" yaml:"max_batch_size"`
	MergeStrategy MergeStrategy `json:"merge_strategy" yaml:"merge_strategy"`
	TargetTopic   string        `json:"target_topic" yaml:"target_topic"`
	Conditions    []Condition   `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	Enabled       bool          `json:"enabled" yaml:"enabled"`
}

type PolicySet struct {
	ByTable map[string]Policy
}
