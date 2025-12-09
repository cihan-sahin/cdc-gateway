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
}

func (r *Router) Route(evt cdcmodel.CDCEvent) (Routed, bool) {
	pol, ok := r.policyStore.GetPolicyForTable(evt.Table)
	if !ok || !pol.Enabled {
		return Routed{}, false
	}
	if !pol.MatchesPayload(evt.Payload) {
		return Routed{}, false
	}
	return Routed{
		Event:  evt,
		Policy: pol,
	}, true
}
