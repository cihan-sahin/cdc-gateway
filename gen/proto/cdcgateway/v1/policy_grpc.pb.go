package cdcv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)


const _ = grpc.SupportPackageIsVersion9

const (
	PolicyService_ListPolicies_FullMethodName = "/cdcgateway.v1.PolicyService/ListPolicies"
	PolicyService_UpsertPolicy_FullMethodName = "/cdcgateway.v1.PolicyService/UpsertPolicy"
	PolicyService_DeletePolicy_FullMethodName = "/cdcgateway.v1.PolicyService/DeletePolicy"
)


type PolicyServiceClient interface {
	ListPolicies(ctx context.Context, in *ListPoliciesRequest, opts ...grpc.CallOption) (*ListPoliciesResponse, error)
	UpsertPolicy(ctx context.Context, in *UpsertPolicyRequest, opts ...grpc.CallOption) (*UpsertPolicyResponse, error)
	DeletePolicy(ctx context.Context, in *DeletePolicyRequest, opts ...grpc.CallOption) (*DeletePolicyResponse, error)
}

type policyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPolicyServiceClient(cc grpc.ClientConnInterface) PolicyServiceClient {
	return &policyServiceClient{cc}
}

func (c *policyServiceClient) ListPolicies(ctx context.Context, in *ListPoliciesRequest, opts ...grpc.CallOption) (*ListPoliciesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPoliciesResponse)
	err := c.cc.Invoke(ctx, PolicyService_ListPolicies_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) UpsertPolicy(ctx context.Context, in *UpsertPolicyRequest, opts ...grpc.CallOption) (*UpsertPolicyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpsertPolicyResponse)
	err := c.cc.Invoke(ctx, PolicyService_UpsertPolicy_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *policyServiceClient) DeletePolicy(ctx context.Context, in *DeletePolicyRequest, opts ...grpc.CallOption) (*DeletePolicyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeletePolicyResponse)
	err := c.cc.Invoke(ctx, PolicyService_DeletePolicy_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type PolicyServiceServer interface {
	ListPolicies(context.Context, *ListPoliciesRequest) (*ListPoliciesResponse, error)
	UpsertPolicy(context.Context, *UpsertPolicyRequest) (*UpsertPolicyResponse, error)
	DeletePolicy(context.Context, *DeletePolicyRequest) (*DeletePolicyResponse, error)
	mustEmbedUnimplementedPolicyServiceServer()
}

type UnimplementedPolicyServiceServer struct{}

func (UnimplementedPolicyServiceServer) ListPolicies(context.Context, *ListPoliciesRequest) (*ListPoliciesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method ListPolicies not implemented")
}
func (UnimplementedPolicyServiceServer) UpsertPolicy(context.Context, *UpsertPolicyRequest) (*UpsertPolicyResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UpsertPolicy not implemented")
}
func (UnimplementedPolicyServiceServer) DeletePolicy(context.Context, *DeletePolicyRequest) (*DeletePolicyResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method DeletePolicy not implemented")
}
func (UnimplementedPolicyServiceServer) mustEmbedUnimplementedPolicyServiceServer() {}
func (UnimplementedPolicyServiceServer) testEmbeddedByValue()                       {}

type UnsafePolicyServiceServer interface {
	mustEmbedUnimplementedPolicyServiceServer()
}

func RegisterPolicyServiceServer(s grpc.ServiceRegistrar, srv PolicyServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PolicyService_ServiceDesc, srv)
}

func _PolicyService_ListPolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPoliciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).ListPolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_ListPolicies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).ListPolicies(ctx, req.(*ListPoliciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_UpsertPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpsertPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).UpsertPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_UpsertPolicy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).UpsertPolicy(ctx, req.(*UpsertPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PolicyService_DeletePolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PolicyServiceServer).DeletePolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PolicyService_DeletePolicy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PolicyServiceServer).DeletePolicy(ctx, req.(*DeletePolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var PolicyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cdcgateway.v1.PolicyService",
	HandlerType: (*PolicyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListPolicies",
			Handler:    _PolicyService_ListPolicies_Handler,
		},
		{
			MethodName: "UpsertPolicy",
			Handler:    _PolicyService_UpsertPolicy_Handler,
		},
		{
			MethodName: "DeletePolicy",
			Handler:    _PolicyService_DeletePolicy_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/cdcgateway/v1/policy.proto",
}
