// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: proto/engine_service.proto

package engine_service

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
	EngineService_GetProgramResult_FullMethodName = "/EngineService/GetProgramResult"
	EngineService_GetTestedResult_FullMethodName  = "/EngineService/GetTestedResult"
)

// EngineServiceClient is the client API for EngineService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EngineServiceClient interface {
	// Obtains the result of executing a program in the engine container.
	GetProgramResult(ctx context.Context, in *Program, opts ...grpc.CallOption) (*Result, error)
	// Obtains the result of executing a program against test cases.
	GetTestedResult(ctx context.Context, in *ProgramWithTests, opts ...grpc.CallOption) (*TestResult, error)
}

type engineServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEngineServiceClient(cc grpc.ClientConnInterface) EngineServiceClient {
	return &engineServiceClient{cc}
}

func (c *engineServiceClient) GetProgramResult(ctx context.Context, in *Program, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, EngineService_GetProgramResult_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *engineServiceClient) GetTestedResult(ctx context.Context, in *ProgramWithTests, opts ...grpc.CallOption) (*TestResult, error) {
	out := new(TestResult)
	err := c.cc.Invoke(ctx, EngineService_GetTestedResult_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EngineServiceServer is the server API for EngineService service.
// All implementations must embed UnimplementedEngineServiceServer
// for forward compatibility
type EngineServiceServer interface {
	// Obtains the result of executing a program in the engine container.
	GetProgramResult(context.Context, *Program) (*Result, error)
	// Obtains the result of executing a program against test cases.
	GetTestedResult(context.Context, *ProgramWithTests) (*TestResult, error)
	mustEmbedUnimplementedEngineServiceServer()
}

// UnimplementedEngineServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEngineServiceServer struct {
}

func (UnimplementedEngineServiceServer) GetProgramResult(context.Context, *Program) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProgramResult not implemented")
}
func (UnimplementedEngineServiceServer) GetTestedResult(context.Context, *ProgramWithTests) (*TestResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTestedResult not implemented")
}
func (UnimplementedEngineServiceServer) mustEmbedUnimplementedEngineServiceServer() {}

// UnsafeEngineServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EngineServiceServer will
// result in compilation errors.
type UnsafeEngineServiceServer interface {
	mustEmbedUnimplementedEngineServiceServer()
}

func RegisterEngineServiceServer(s grpc.ServiceRegistrar, srv EngineServiceServer) {
	s.RegisterService(&EngineService_ServiceDesc, srv)
}

func _EngineService_GetProgramResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Program)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EngineServiceServer).GetProgramResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EngineService_GetProgramResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EngineServiceServer).GetProgramResult(ctx, req.(*Program))
	}
	return interceptor(ctx, in, info, handler)
}

func _EngineService_GetTestedResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProgramWithTests)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EngineServiceServer).GetTestedResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EngineService_GetTestedResult_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EngineServiceServer).GetTestedResult(ctx, req.(*ProgramWithTests))
	}
	return interceptor(ctx, in, info, handler)
}

// EngineService_ServiceDesc is the grpc.ServiceDesc for EngineService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EngineService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "EngineService",
	HandlerType: (*EngineServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProgramResult",
			Handler:    _EngineService_GetProgramResult_Handler,
		},
		{
			MethodName: "GetTestedResult",
			Handler:    _EngineService_GetTestedResult_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/engine_service.proto",
}