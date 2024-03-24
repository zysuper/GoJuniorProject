// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: reward/v1/reward.proto

package rewardv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	RewardService_PreReward_FullMethodName = "/reward.v1.RewardService/PreReward"
	RewardService_GetReward_FullMethodName = "/reward.v1.RewardService/GetReward"
)

// RewardServiceClient is the client API for RewardService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RewardServiceClient interface {
	PreReward(ctx context.Context, in *PreRewardRequest, opts ...grpc.CallOption) (*PreRewardResponse, error)
	GetReward(ctx context.Context, in *GetRewardRequest, opts ...grpc.CallOption) (*GetRewardResponse, error)
}

type rewardServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRewardServiceClient(cc grpc.ClientConnInterface) RewardServiceClient {
	return &rewardServiceClient{cc}
}

func (c *rewardServiceClient) PreReward(ctx context.Context, in *PreRewardRequest, opts ...grpc.CallOption) (*PreRewardResponse, error) {
	out := new(PreRewardResponse)
	err := c.cc.Invoke(ctx, RewardService_PreReward_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rewardServiceClient) GetReward(ctx context.Context, in *GetRewardRequest, opts ...grpc.CallOption) (*GetRewardResponse, error) {
	out := new(GetRewardResponse)
	err := c.cc.Invoke(ctx, RewardService_GetReward_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RewardServiceServer is the server API for RewardService service.
// All implementations must embed UnimplementedRewardServiceServer
// for forward compatibility
type RewardServiceServer interface {
	PreReward(context.Context, *PreRewardRequest) (*PreRewardResponse, error)
	GetReward(context.Context, *GetRewardRequest) (*GetRewardResponse, error)
	mustEmbedUnimplementedRewardServiceServer()
}

// UnimplementedRewardServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRewardServiceServer struct {
}

func (UnimplementedRewardServiceServer) PreReward(context.Context, *PreRewardRequest) (*PreRewardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreReward not implemented")
}
func (UnimplementedRewardServiceServer) GetReward(context.Context, *GetRewardRequest) (*GetRewardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReward not implemented")
}
func (UnimplementedRewardServiceServer) mustEmbedUnimplementedRewardServiceServer() {}

// UnsafeRewardServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RewardServiceServer will
// result in compilation errors.
type UnsafeRewardServiceServer interface {
	mustEmbedUnimplementedRewardServiceServer()
}

func RegisterRewardServiceServer(s grpc.ServiceRegistrar, srv RewardServiceServer) {
	s.RegisterService(&RewardService_ServiceDesc, srv)
}

func _RewardService_PreReward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PreRewardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RewardServiceServer).PreReward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RewardService_PreReward_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RewardServiceServer).PreReward(ctx, req.(*PreRewardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RewardService_GetReward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRewardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RewardServiceServer).GetReward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RewardService_GetReward_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RewardServiceServer).GetReward(ctx, req.(*GetRewardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RewardService_ServiceDesc is the grpc.ServiceDesc for RewardService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RewardService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "reward.v1.RewardService",
	HandlerType: (*RewardServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PreReward",
			Handler:    _RewardService_PreReward_Handler,
		},
		{
			MethodName: "GetReward",
			Handler:    _RewardService_GetReward_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "reward/v1/reward.proto",
}