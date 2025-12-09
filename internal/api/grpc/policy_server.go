package grpcapi

import (
	"context"

	cdcv1 "github.com/cihan-sahin/cdc-gateway/gen/cdcgateway/v1"
	"github.com/cihan-sahin/cdc-gateway/internal/policy"
)

type PolicyServer struct {
	cdcv1.UnimplementedPolicyServiceServer
	store policy.Store
}

func NewPolicyServer(store policy.Store) *PolicyServer {
	return &PolicyServer{store: store}
}

func (s *PolicyServer) ListPolicies(ctx context.Context, req *cdcv1.ListPoliciesRequest) (*cdcv1.ListPoliciesResponse, error) {
	pols := s.store.All()

	out := make([]*cdcv1.Policy, 0, len(pols))
	for _, p := range pols {
		out = append(out, &cdcv1.Policy{
			Table:         p.Table,
			Mode:          string(p.Mode),
			WindowMs:      int32(p.WindowMs),
			MaxBatchSize:  int32(p.MaxBatchSize),
			MergeStrategy: string(p.MergeStrategy),
			TargetTopic:   p.TargetTopic,
			Enabled:       p.Enabled,
		})
	}

	return &cdcv1.ListPoliciesResponse{
		Policies: out,
	}, nil
}

func (s *PolicyServer) UpsertPolicy(ctx context.Context, req *cdcv1.UpsertPolicyRequest) (*cdcv1.UpsertPolicyResponse, error) {
	in := req.GetPolicy()
	if in == nil {
		return &cdcv1.UpsertPolicyResponse{}, nil
	}

	p := policy.Policy{
		Table:         in.GetTable(),
		Mode:          policy.Mode(in.GetMode()),
		WindowMs:      int(in.GetWindowMs()),
		MaxBatchSize:  int(in.GetMaxBatchSize()),
		MergeStrategy: policy.MergeStrategy(in.GetMergeStrategy()),
		TargetTopic:   in.GetTargetTopic(),
		Enabled:       in.GetEnabled(),
	}

	if err := s.store.Upsert(p); err != nil {
		return nil, err
	}

	return &cdcv1.UpsertPolicyResponse{}, nil
}

func (s *PolicyServer) DeletePolicy(ctx context.Context, req *cdcv1.DeletePolicyRequest) (*cdcv1.DeletePolicyResponse, error) {
	table := req.GetTable()
	if table != "" {
		_ = s.store.Delete(table)
	}
	return &cdcv1.DeletePolicyResponse{}, nil
}
