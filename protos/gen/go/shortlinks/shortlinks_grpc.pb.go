// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: shortlinks/shortlinks.proto

package shortlinksv1

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

// ShortlinksClient is the client API for Shortlinks service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortlinksClient interface {
	GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error)
	SaveURL(ctx context.Context, in *SaveURLRequest, opts ...grpc.CallOption) (*SaveURLResponse, error)
}

type shortlinksClient struct {
	cc grpc.ClientConnInterface
}

func NewShortlinksClient(cc grpc.ClientConnInterface) ShortlinksClient {
	return &shortlinksClient{cc}
}

func (c *shortlinksClient) GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error) {
	out := new(GetURLResponse)
	err := c.cc.Invoke(ctx, "/shortlinks.Shortlinks/GetURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortlinksClient) SaveURL(ctx context.Context, in *SaveURLRequest, opts ...grpc.CallOption) (*SaveURLResponse, error) {
	out := new(SaveURLResponse)
	err := c.cc.Invoke(ctx, "/shortlinks.Shortlinks/SaveURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortlinksServer is the server API for Shortlinks service.
// All implementations must embed UnimplementedShortlinksServer
// for forward compatibility
type ShortlinksServer interface {
	GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error)
	SaveURL(context.Context, *SaveURLRequest) (*SaveURLResponse, error)
	mustEmbedUnimplementedShortlinksServer()
}

// UnimplementedShortlinksServer must be embedded to have forward compatible implementations.
type UnimplementedShortlinksServer struct {
}

func (UnimplementedShortlinksServer) GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURL not implemented")
}
func (UnimplementedShortlinksServer) SaveURL(context.Context, *SaveURLRequest) (*SaveURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveURL not implemented")
}
func (UnimplementedShortlinksServer) mustEmbedUnimplementedShortlinksServer() {}

// UnsafeShortlinksServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortlinksServer will
// result in compilation errors.
type UnsafeShortlinksServer interface {
	mustEmbedUnimplementedShortlinksServer()
}

func RegisterShortlinksServer(s grpc.ServiceRegistrar, srv ShortlinksServer) {
	s.RegisterService(&Shortlinks_ServiceDesc, srv)
}

func _Shortlinks_GetURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortlinksServer).GetURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortlinks.Shortlinks/GetURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortlinksServer).GetURL(ctx, req.(*GetURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Shortlinks_SaveURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortlinksServer).SaveURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortlinks.Shortlinks/SaveURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortlinksServer).SaveURL(ctx, req.(*SaveURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Shortlinks_ServiceDesc is the grpc.ServiceDesc for Shortlinks service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Shortlinks_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortlinks.Shortlinks",
	HandlerType: (*ShortlinksServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetURL",
			Handler:    _Shortlinks_GetURL_Handler,
		},
		{
			MethodName: "SaveURL",
			Handler:    _Shortlinks_SaveURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shortlinks/shortlinks.proto",
}
