package pipeline

import (
	"time"

	"github.com/cihan-sahin/cdc-gateway/internal/policy"
	"github.com/cihan-sahin/cdc-gateway/pkg/cdcmodel"
)

type CoalescedResult struct {
	Events       []cdcmodel.CDCEvent
	Policy       policy.Policy
	Partition    int32
	FlushedAt    time.Time
	WindowOpened time.Time
}

type Coalescer interface {
	AddEvent(evt cdcmodel.CDCEvent, pol policy.Policy, partition int32, now time.Time) ([]CoalescedResult, error)

	FlushDue(now time.Time) ([]CoalescedResult, error)
}
