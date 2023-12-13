// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// DiagnosticClient is the client API for Diagnostic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DiagnosticClient interface {
	Screenshot(ctx context.Context, in *ScreenshotRequest, opts ...grpc.CallOption) (*ScreenshotResponse, error)
	ClearSystemEventLog(ctx context.Context, in *ClearSystemEventLogRequest, opts ...grpc.CallOption) (*ClearSystemEventLogResponse, error)
	SystemEventLog(ctx context.Context, in *SystemEventLogRequest, opts ...grpc.CallOption) (*SystemEventLogResponse, error)
	SystemEventLogRaw(ctx context.Context, in *SystemEventLogRawRequest, opts ...grpc.CallOption) (*SystemEventLogRawResponse, error)
}

type diagnosticClient struct {
	cc grpc.ClientConnInterface
}

func NewDiagnosticClient(cc grpc.ClientConnInterface) DiagnosticClient {
	return &diagnosticClient{cc}
}

func (c *diagnosticClient) Screenshot(ctx context.Context, in *ScreenshotRequest, opts ...grpc.CallOption) (*ScreenshotResponse, error) {
	out := new(ScreenshotResponse)
	err := c.cc.Invoke(ctx, "/github.com.tinkerbell.pbnj.api.v1.Diagnostic/Screenshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *diagnosticClient) ClearSystemEventLog(ctx context.Context, in *ClearSystemEventLogRequest, opts ...grpc.CallOption) (*ClearSystemEventLogResponse, error) {
	out := new(ClearSystemEventLogResponse)
	err := c.cc.Invoke(ctx, "/github.com.tinkerbell.pbnj.api.v1.Diagnostic/ClearSystemEventLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *diagnosticClient) SystemEventLog(ctx context.Context, in *SystemEventLogRequest, opts ...grpc.CallOption) (*SystemEventLogResponse, error) {
	out := new(SystemEventLogResponse)
	err := c.cc.Invoke(ctx, "/github.com.tinkerbell.pbnj.api.v1.Diagnostic/SystemEventLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *diagnosticClient) SystemEventLogRaw(ctx context.Context, in *SystemEventLogRawRequest, opts ...grpc.CallOption) (*SystemEventLogRawResponse, error) {
	out := new(SystemEventLogRawResponse)
	err := c.cc.Invoke(ctx, "/github.com.tinkerbell.pbnj.api.v1.Diagnostic/SystemEventLogRaw", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiagnosticServer is the server API for Diagnostic service.
// All implementations must embed UnimplementedDiagnosticServer
// for forward compatibility
type DiagnosticServer interface {
	Screenshot(context.Context, *ScreenshotRequest) (*ScreenshotResponse, error)
	ClearSystemEventLog(context.Context, *ClearSystemEventLogRequest) (*ClearSystemEventLogResponse, error)
	SystemEventLog(context.Context, *SystemEventLogRequest) (*SystemEventLogResponse, error)
	SystemEventLogRaw(context.Context, *SystemEventLogRawRequest) (*SystemEventLogRawResponse, error)
	mustEmbedUnimplementedDiagnosticServer()
}

// UnimplementedDiagnosticServer must be embedded to have forward compatible implementations.
type UnimplementedDiagnosticServer struct {
}

func (UnimplementedDiagnosticServer) Screenshot(context.Context, *ScreenshotRequest) (*ScreenshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Screenshot not implemented")
}
func (UnimplementedDiagnosticServer) ClearSystemEventLog(context.Context, *ClearSystemEventLogRequest) (*ClearSystemEventLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearSystemEventLog not implemented")
}
func (UnimplementedDiagnosticServer) SystemEventLog(context.Context, *SystemEventLogRequest) (*SystemEventLogResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SystemEventLog not implemented")
}
func (UnimplementedDiagnosticServer) SystemEventLogRaw(context.Context, *SystemEventLogRawRequest) (*SystemEventLogRawResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SystemEventLogRaw not implemented")
}
func (UnimplementedDiagnosticServer) mustEmbedUnimplementedDiagnosticServer() {}

// UnsafeDiagnosticServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DiagnosticServer will
// result in compilation errors.
type UnsafeDiagnosticServer interface {
	mustEmbedUnimplementedDiagnosticServer()
}

func RegisterDiagnosticServer(s grpc.ServiceRegistrar, srv DiagnosticServer) {
	s.RegisterService(&_Diagnostic_serviceDesc, srv)
}

func _Diagnostic_Screenshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScreenshotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiagnosticServer).Screenshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/github.com.tinkerbell.pbnj.api.v1.Diagnostic/Screenshot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiagnosticServer).Screenshot(ctx, req.(*ScreenshotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Diagnostic_ClearSystemEventLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearSystemEventLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiagnosticServer).ClearSystemEventLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/github.com.tinkerbell.pbnj.api.v1.Diagnostic/ClearSystemEventLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiagnosticServer).ClearSystemEventLog(ctx, req.(*ClearSystemEventLogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Diagnostic_SystemEventLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SystemEventLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiagnosticServer).SystemEventLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/github.com.tinkerbell.pbnj.api.v1.Diagnostic/SystemEventLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiagnosticServer).SystemEventLog(ctx, req.(*SystemEventLogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Diagnostic_SystemEventLogRaw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SystemEventLogRawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiagnosticServer).SystemEventLogRaw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/github.com.tinkerbell.pbnj.api.v1.Diagnostic/SystemEventLogRaw",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiagnosticServer).SystemEventLogRaw(ctx, req.(*SystemEventLogRawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Diagnostic_serviceDesc = grpc.ServiceDesc{
	ServiceName: "github.com.tinkerbell.pbnj.api.v1.Diagnostic",
	HandlerType: (*DiagnosticServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Screenshot",
			Handler:    _Diagnostic_Screenshot_Handler,
		},
		{
			MethodName: "ClearSystemEventLog",
			Handler:    _Diagnostic_ClearSystemEventLog_Handler,
		},
		{
			MethodName: "SystemEventLog",
			Handler:    _Diagnostic_SystemEventLog_Handler,
		},
		{
			MethodName: "SystemEventLogRaw",
			Handler:    _Diagnostic_SystemEventLogRaw_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/diagnostic.proto",
}
