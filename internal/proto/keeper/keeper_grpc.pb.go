// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// KeeperClient is the client API for Keeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeeperClient interface {
	Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type keeperClient struct {
	cc grpc.ClientConnInterface
}

func NewKeeperClient(cc grpc.ClientConnInterface) KeeperClient {
	return &keeperClient{cc}
}

func (c *keeperClient) Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchResponse, error) {
	out := new(FetchResponse)
	err := c.cc.Invoke(ctx, "/proto.Keeper/Fetch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.Keeper/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeeperServer is the server API for Keeper service.
// All implementations must embed UnimplementedKeeperServer
// for forward compatibility
type KeeperServer interface {
	Fetch(context.Context, *FetchRequest) (*FetchResponse, error)
	Create(context.Context, *CreateRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedKeeperServer()
}

// UnimplementedKeeperServer must be embedded to have forward compatible implementations.
type UnimplementedKeeperServer struct {
}

func (UnimplementedKeeperServer) Fetch(context.Context, *FetchRequest) (*FetchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fetch not implemented")
}
func (UnimplementedKeeperServer) Create(context.Context, *CreateRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedKeeperServer) mustEmbedUnimplementedKeeperServer() {}

// UnsafeKeeperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeeperServer will
// result in compilation errors.
type UnsafeKeeperServer interface {
	mustEmbedUnimplementedKeeperServer()
}

func RegisterKeeperServer(s grpc.ServiceRegistrar, srv KeeperServer) {
	s.RegisterService(&Keeper_ServiceDesc, srv)
}

func _Keeper_Fetch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).Fetch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Keeper/Fetch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).Fetch(ctx, req.(*FetchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Keeper/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Keeper_ServiceDesc is the grpc.ServiceDesc for Keeper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Keeper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Keeper",
	HandlerType: (*KeeperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Fetch",
			Handler:    _Keeper_Fetch_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Keeper_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/keeper/keeper.proto",
}
