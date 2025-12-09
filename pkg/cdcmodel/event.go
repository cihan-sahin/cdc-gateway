package cdcmodel

type Operation string

const (
	OpCreate Operation = "c"
	OpUpdate Operation = "u"
	OpDelete Operation = "d"
	OpRead   Operation = "r"
)

type CDCEvent struct {
	Key string

	Table string

	Op Operation

	SourceTopic string

	Payload map[string]interface{}
}
