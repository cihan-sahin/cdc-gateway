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

// Her partition için bir coalescer instance'ı olacak.
type Coalescer interface {
	// Yeni event geldiğinde çağrılır.
	// Gerekirse sadece buffer'a koyar; flush zamanı gelmişse
	// dönen slice'te flush edilmesi gereken eventler olur.
	AddEvent(evt cdcmodel.CDCEvent, pol policy.Policy, partition int32, now time.Time) ([]CoalescedResult, error)

	// Zaman bazlı flush için (örn. her 100ms tick).
	FlushDue(now time.Time) ([]CoalescedResult, error)
}
