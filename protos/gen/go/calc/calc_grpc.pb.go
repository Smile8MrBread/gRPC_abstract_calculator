// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: calc.proto

package calcv1

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
	Calc_SaveExpression_FullMethodName         = "/calc.Calc/SaveExpression"
	Calc_GetExpression_FullMethodName          = "/calc.Calc/GetExpression"
	Calc_UpdateExpression_FullMethodName       = "/calc.Calc/UpdateExpression"
	Calc_DeleteExpression_FullMethodName       = "/calc.Calc/DeleteExpression"
	Calc_GetNotDonedExpressions_FullMethodName = "/calc.Calc/GetNotDonedExpressions"
	Calc_GetAllExpressions_FullMethodName      = "/calc.Calc/GetAllExpressions"
	Calc_UpdateArithmetic_FullMethodName       = "/calc.Calc/UpdateArithmetic"
	Calc_GetArithmetic_FullMethodName          = "/calc.Calc/GetArithmetic"
)

// CalcClient is the client API for Calc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalcClient interface {
	SaveExpression(ctx context.Context, in *SaveExpressionRequest, opts ...grpc.CallOption) (*SaveExpressionResponse, error)
	GetExpression(ctx context.Context, in *GetExpressionRequest, opts ...grpc.CallOption) (*GetExpressionResponse, error)
	UpdateExpression(ctx context.Context, in *UpdateExpressionRequest, opts ...grpc.CallOption) (*NothingMessage, error)
	DeleteExpression(ctx context.Context, in *DeleteExpressionRequest, opts ...grpc.CallOption) (*NothingMessage, error)
	GetNotDonedExpressions(ctx context.Context, in *NothingMessage, opts ...grpc.CallOption) (*GetNotDonedExpressionsResponse, error)
	GetAllExpressions(ctx context.Context, in *NothingMessage, opts ...grpc.CallOption) (*GetAllExpressionsResponse, error)
	UpdateArithmetic(ctx context.Context, in *UpdateArithmeticRequest, opts ...grpc.CallOption) (*NothingMessage, error)
	GetArithmetic(ctx context.Context, in *GetArithmeticRequest, opts ...grpc.CallOption) (*GetArithmeticResponse, error)
}

type calcClient struct {
	cc grpc.ClientConnInterface
}

func NewCalcClient(cc grpc.ClientConnInterface) CalcClient {
	return &calcClient{cc}
}

func (c *calcClient) SaveExpression(ctx context.Context, in *SaveExpressionRequest, opts ...grpc.CallOption) (*SaveExpressionResponse, error) {
	out := new(SaveExpressionResponse)
	err := c.cc.Invoke(ctx, Calc_SaveExpression_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calcClient) GetExpression(ctx context.Context, in *GetExpressionRequest, opts ...grpc.CallOption) (*GetExpressionResponse, error) {
	out := new(GetExpressionResponse)
	err := c.cc.Invoke(ctx, Calc_GetExpression_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calcClient) UpdateExpression(ctx context.Context, in *UpdateExpressionRequest, opts ...grpc.CallOption) (*NothingMessage, error) {
	out := new(NothingMessage)
	err := c.cc.Invoke(ctx, Calc_UpdateExpression_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calcClient) DeleteExpression(ctx context.Context, in *DeleteExpressionRequest, opts ...grpc.CallOption) (*NothingMessage, error) {
	out := new(NothingMessage)
	err := c.cc.Invoke(ctx, Calc_DeleteExpression_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calcClient) GetNotDonedExpressions(ctx context.Context, in *NothingMessage, opts ...grpc.CallOption) (*GetNotDonedExpressionsResponse, error) {
	out := new(GetNotDonedExpressionsResponse)
	err := c.cc.Invoke(ctx, Calc_GetNotDonedExpressions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calcClient) GetAllExpressions(ctx context.Context, in *NothingMessage, opts ...grpc.CallOption) (*GetAllExpressionsResponse, error) {
	out := new(GetAllExpressionsResponse)
	err := c.cc.Invoke(ctx, Calc_GetAllExpressions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calcClient) UpdateArithmetic(ctx context.Context, in *UpdateArithmeticRequest, opts ...grpc.CallOption) (*NothingMessage, error) {
	out := new(NothingMessage)
	err := c.cc.Invoke(ctx, Calc_UpdateArithmetic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calcClient) GetArithmetic(ctx context.Context, in *GetArithmeticRequest, opts ...grpc.CallOption) (*GetArithmeticResponse, error) {
	out := new(GetArithmeticResponse)
	err := c.cc.Invoke(ctx, Calc_GetArithmetic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalcServer is the server API for Calc service.
// All implementations must embed UnimplementedCalcServer
// for forward compatibility
type CalcServer interface {
	SaveExpression(context.Context, *SaveExpressionRequest) (*SaveExpressionResponse, error)
	GetExpression(context.Context, *GetExpressionRequest) (*GetExpressionResponse, error)
	UpdateExpression(context.Context, *UpdateExpressionRequest) (*NothingMessage, error)
	DeleteExpression(context.Context, *DeleteExpressionRequest) (*NothingMessage, error)
	GetNotDonedExpressions(context.Context, *NothingMessage) (*GetNotDonedExpressionsResponse, error)
	GetAllExpressions(context.Context, *NothingMessage) (*GetAllExpressionsResponse, error)
	UpdateArithmetic(context.Context, *UpdateArithmeticRequest) (*NothingMessage, error)
	GetArithmetic(context.Context, *GetArithmeticRequest) (*GetArithmeticResponse, error)
	mustEmbedUnimplementedCalcServer()
}

// UnimplementedCalcServer must be embedded to have forward compatible implementations.
type UnimplementedCalcServer struct {
}

func (UnimplementedCalcServer) SaveExpression(context.Context, *SaveExpressionRequest) (*SaveExpressionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveExpression not implemented")
}
func (UnimplementedCalcServer) GetExpression(context.Context, *GetExpressionRequest) (*GetExpressionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExpression not implemented")
}
func (UnimplementedCalcServer) UpdateExpression(context.Context, *UpdateExpressionRequest) (*NothingMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateExpression not implemented")
}
func (UnimplementedCalcServer) DeleteExpression(context.Context, *DeleteExpressionRequest) (*NothingMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteExpression not implemented")
}
func (UnimplementedCalcServer) GetNotDonedExpressions(context.Context, *NothingMessage) (*GetNotDonedExpressionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotDonedExpressions not implemented")
}
func (UnimplementedCalcServer) GetAllExpressions(context.Context, *NothingMessage) (*GetAllExpressionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllExpressions not implemented")
}
func (UnimplementedCalcServer) UpdateArithmetic(context.Context, *UpdateArithmeticRequest) (*NothingMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateArithmetic not implemented")
}
func (UnimplementedCalcServer) GetArithmetic(context.Context, *GetArithmeticRequest) (*GetArithmeticResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArithmetic not implemented")
}
func (UnimplementedCalcServer) mustEmbedUnimplementedCalcServer() {}

// UnsafeCalcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalcServer will
// result in compilation errors.
type UnsafeCalcServer interface {
	mustEmbedUnimplementedCalcServer()
}

func RegisterCalcServer(s grpc.ServiceRegistrar, srv CalcServer) {
	s.RegisterService(&Calc_ServiceDesc, srv)
}

func _Calc_SaveExpression_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveExpressionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcServer).SaveExpression(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calc_SaveExpression_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcServer).SaveExpression(ctx, req.(*SaveExpressionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calc_GetExpression_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetExpressionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcServer).GetExpression(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calc_GetExpression_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcServer).GetExpression(ctx, req.(*GetExpressionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calc_UpdateExpression_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateExpressionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcServer).UpdateExpression(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calc_UpdateExpression_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcServer).UpdateExpression(ctx, req.(*UpdateExpressionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calc_DeleteExpression_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteExpressionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcServer).DeleteExpression(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calc_DeleteExpression_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcServer).DeleteExpression(ctx, req.(*DeleteExpressionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calc_GetNotDonedExpressions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NothingMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcServer).GetNotDonedExpressions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calc_GetNotDonedExpressions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcServer).GetNotDonedExpressions(ctx, req.(*NothingMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calc_GetAllExpressions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NothingMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcServer).GetAllExpressions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calc_GetAllExpressions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcServer).GetAllExpressions(ctx, req.(*NothingMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calc_UpdateArithmetic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateArithmeticRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcServer).UpdateArithmetic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calc_UpdateArithmetic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcServer).UpdateArithmetic(ctx, req.(*UpdateArithmeticRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Calc_GetArithmetic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetArithmeticRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalcServer).GetArithmetic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Calc_GetArithmetic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalcServer).GetArithmetic(ctx, req.(*GetArithmeticRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Calc_ServiceDesc is the grpc.ServiceDesc for Calc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Calc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "calc.Calc",
	HandlerType: (*CalcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveExpression",
			Handler:    _Calc_SaveExpression_Handler,
		},
		{
			MethodName: "GetExpression",
			Handler:    _Calc_GetExpression_Handler,
		},
		{
			MethodName: "UpdateExpression",
			Handler:    _Calc_UpdateExpression_Handler,
		},
		{
			MethodName: "DeleteExpression",
			Handler:    _Calc_DeleteExpression_Handler,
		},
		{
			MethodName: "GetNotDonedExpressions",
			Handler:    _Calc_GetNotDonedExpressions_Handler,
		},
		{
			MethodName: "GetAllExpressions",
			Handler:    _Calc_GetAllExpressions_Handler,
		},
		{
			MethodName: "UpdateArithmetic",
			Handler:    _Calc_UpdateArithmetic_Handler,
		},
		{
			MethodName: "GetArithmetic",
			Handler:    _Calc_GetArithmetic_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "calc.proto",
}
