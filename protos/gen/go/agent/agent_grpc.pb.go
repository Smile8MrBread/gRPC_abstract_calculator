// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: agent.proto

package agentv1

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
	GRPCAgent_ExpForDo_FullMethodName = "/agent.GRPCAgent/ExpForDo"
)

// GRPCAgentClient is the client API for GRPCAgent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GRPCAgentClient interface {
	ExpForDo(ctx context.Context, in *ExpForDoRequest, opts ...grpc.CallOption) (*NothingMessage, error)
}

type gRPCAgentClient struct {
	cc grpc.ClientConnInterface
}

func NewGRPCAgentClient(cc grpc.ClientConnInterface) GRPCAgentClient {
	return &gRPCAgentClient{cc}
}

func (c *gRPCAgentClient) ExpForDo(ctx context.Context, in *ExpForDoRequest, opts ...grpc.CallOption) (*NothingMessage, error) {
	out := new(NothingMessage)
	err := c.cc.Invoke(ctx, GRPCAgent_ExpForDo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GRPCAgentServer is the server API for GRPCAgent service.
// All implementations must embed UnimplementedGRPCAgentServer
// for forward compatibility
type GRPCAgentServer interface {
	ExpForDo(context.Context, *ExpForDoRequest) (*NothingMessage, error)
	mustEmbedUnimplementedGRPCAgentServer()
}

// UnimplementedGRPCAgentServer must be embedded to have forward compatible implementations.
type UnimplementedGRPCAgentServer struct {
}

func (UnimplementedGRPCAgentServer) ExpForDo(context.Context, *ExpForDoRequest) (*NothingMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExpForDo not implemented")
}
func (UnimplementedGRPCAgentServer) mustEmbedUnimplementedGRPCAgentServer() {}

// UnsafeGRPCAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GRPCAgentServer will
// result in compilation errors.
type UnsafeGRPCAgentServer interface {
	mustEmbedUnimplementedGRPCAgentServer()
}

func RegisterGRPCAgentServer(s grpc.ServiceRegistrar, srv GRPCAgentServer) {
	s.RegisterService(&GRPCAgent_ServiceDesc, srv)
}

func _GRPCAgent_ExpForDo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExpForDoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GRPCAgentServer).ExpForDo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GRPCAgent_ExpForDo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GRPCAgentServer).ExpForDo(ctx, req.(*ExpForDoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GRPCAgent_ServiceDesc is the grpc.ServiceDesc for GRPCAgent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GRPCAgent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "agent.GRPCAgent",
	HandlerType: (*GRPCAgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExpForDo",
			Handler:    _GRPCAgent_ExpForDo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "agent.proto",
}
