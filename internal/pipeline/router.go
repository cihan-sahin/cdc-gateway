package pipeline

import (
	"github.com/cihan-sahin/cdc-gateway/internal/policy"
	"github.com/cihan-sahin/cdc-gateway/pkg/cdcmodel"
)

type Router struct {
	policyStore policy.Store
}

func NewRouter(store policy.Store) *Router {
	return &Router{policyStore: store}
}

type Routed struct {
	Event  cdcmodel.CDCEvent
	Policy policy.Policy
	// Mode: matched policy mi, yoksa default mu vs. ileride ekleyebiliriz
}

// Route: tabloya göre policy bul, conditions tutuyorsa döndür
// Policy yoksa ok==false dönüyoruz (pipeline karar verecek: passthrough mu drop mu)
func (r *Router) Route(evt cdcmodel.CDCEvent) (Routed, bool) {
	pol, ok := r.policyStore.GetPolicyForTable(evt.Table)
	if !ok || !pol.Enabled {
		return Routed{}, false
	}
	// conditions payload'a uyuyor mu?
	if !pol.MatchesPayload(evt.Payload) {
		return Routed{}, false
	}
	return Routed{
		Event:  evt,
		Policy: pol,
	}, true
}
