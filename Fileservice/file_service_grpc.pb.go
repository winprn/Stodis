// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/file_service.proto

package Fileservice

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UploadFile_CreateFile_FullMethodName = "/fileservice.UploadFile/CreateFile"
	UploadFile_UploadFile_FullMethodName = "/fileservice.UploadFile/UploadFile"
)

// UploadFileClient is the client API for UploadFile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UploadFileClient interface {
	CreateFile(ctx context.Context, in *CreateFileRequest, opts ...grpc.CallOption) (*CreateFileResponse, error)
	// Send a single file
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileData, FileUploadResponse], error)
}

type uploadFileClient struct {
	cc grpc.ClientConnInterface
}

func NewUploadFileClient(cc grpc.ClientConnInterface) UploadFileClient {
	return &uploadFileClient{cc}
}

func (c *uploadFileClient) CreateFile(ctx context.Context, in *CreateFileRequest, opts ...grpc.CallOption) (*CreateFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateFileResponse)
	err := c.cc.Invoke(ctx, UploadFile_CreateFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadFileClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileData, FileUploadResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &UploadFile_ServiceDesc.Streams[0], UploadFile_UploadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[FileData, FileUploadResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type UploadFile_UploadFileClient = grpc.ClientStreamingClient[FileData, FileUploadResponse]

// UploadFileServer is the server API for UploadFile service.
// All implementations must embed UnimplementedUploadFileServer
// for forward compatibility.
type UploadFileServer interface {
	CreateFile(context.Context, *CreateFileRequest) (*CreateFileResponse, error)
	// Send a single file
	UploadFile(grpc.ClientStreamingServer[FileData, FileUploadResponse]) error
	mustEmbedUnimplementedUploadFileServer()
}

// UnimplementedUploadFileServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUploadFileServer struct{}

func (UnimplementedUploadFileServer) CreateFile(context.Context, *CreateFileRequest) (*CreateFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateFile not implemented")
}
func (UnimplementedUploadFileServer) UploadFile(grpc.ClientStreamingServer[FileData, FileUploadResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedUploadFileServer) mustEmbedUnimplementedUploadFileServer() {}
func (UnimplementedUploadFileServer) testEmbeddedByValue()                    {}

// UnsafeUploadFileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UploadFileServer will
// result in compilation errors.
type UnsafeUploadFileServer interface {
	mustEmbedUnimplementedUploadFileServer()
}

func RegisterUploadFileServer(s grpc.ServiceRegistrar, srv UploadFileServer) {
	// If the following call pancis, it indicates UnimplementedUploadFileServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UploadFile_ServiceDesc, srv)
}

func _UploadFile_CreateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadFileServer).CreateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UploadFile_CreateFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadFileServer).CreateFile(ctx, req.(*CreateFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UploadFile_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UploadFileServer).UploadFile(&grpc.GenericServerStream[FileData, FileUploadResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type UploadFile_UploadFileServer = grpc.ClientStreamingServer[FileData, FileUploadResponse]

// UploadFile_ServiceDesc is the grpc.ServiceDesc for UploadFile service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UploadFile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "fileservice.UploadFile",
	HandlerType: (*UploadFileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateFile",
			Handler:    _UploadFile_CreateFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _UploadFile_UploadFile_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "proto/file_service.proto",
}
