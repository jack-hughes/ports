// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package ports

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// PortsClient is the client API for Ports service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PortsClient interface {
	Update(ctx context.Context, opts ...grpc.CallOption) (Ports_UpdateClient, error)
	Get(ctx context.Context, in *GetPortRequest, opts ...grpc.CallOption) (*Port, error)
	List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Ports_ListClient, error)
}

type portsClient struct {
	cc grpc.ClientConnInterface
}

func NewPortsClient(cc grpc.ClientConnInterface) PortsClient {
	return &portsClient{cc}
}

func (c *portsClient) Update(ctx context.Context, opts ...grpc.CallOption) (Ports_UpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Ports_serviceDesc.Streams[0], "/ports.v1.Ports/Update", opts...)
	if err != nil {
		return nil, err
	}
	x := &portsUpdateClient{stream}
	return x, nil
}

type Ports_UpdateClient interface {
	Send(*Port) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type portsUpdateClient struct {
	grpc.ClientStream
}

func (x *portsUpdateClient) Send(m *Port) error {
	return x.ClientStream.SendMsg(m)
}

func (x *portsUpdateClient) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *portsClient) Get(ctx context.Context, in *GetPortRequest, opts ...grpc.CallOption) (*Port, error) {
	out := new(Port)
	err := c.cc.Invoke(ctx, "/ports.v1.Ports/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *portsClient) List(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Ports_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Ports_serviceDesc.Streams[1], "/ports.v1.Ports/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &portsListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Ports_ListClient interface {
	Recv() (*Port, error)
	grpc.ClientStream
}

type portsListClient struct {
	grpc.ClientStream
}

func (x *portsListClient) Recv() (*Port, error) {
	m := new(Port)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PortsServer is the server API for Ports service.
// All implementations must embed UnimplementedPortsServer
// for forward compatibility
type PortsServer interface {
	Update(Ports_UpdateServer) error
	Get(context.Context, *GetPortRequest) (*Port, error)
	List(*emptypb.Empty, Ports_ListServer) error
	mustEmbedUnimplementedPortsServer()
}

// UnimplementedPortsServer must be embedded to have forward compatible implementations.
type UnimplementedPortsServer struct {
}

func (UnimplementedPortsServer) Update(Ports_UpdateServer) error {
	return status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedPortsServer) Get(context.Context, *GetPortRequest) (*Port, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedPortsServer) List(*emptypb.Empty, Ports_ListServer) error {
	return status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedPortsServer) mustEmbedUnimplementedPortsServer() {}

// UnsafePortsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PortsServer will
// result in compilation errors.
type UnsafePortsServer interface {
	mustEmbedUnimplementedPortsServer()
}

func RegisterPortsServer(s *grpc.Server, srv PortsServer) {
	s.RegisterService(&_Ports_serviceDesc, srv)
}

func _Ports_Update_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PortsServer).Update(&portsUpdateServer{stream})
}

type Ports_UpdateServer interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*Port, error)
	grpc.ServerStream
}

type portsUpdateServer struct {
	grpc.ServerStream
}

func (x *portsUpdateServer) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *portsUpdateServer) Recv() (*Port, error) {
	m := new(Port)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Ports_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPortRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PortsServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ports.v1.Ports/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PortsServer).Get(ctx, req.(*GetPortRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ports_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PortsServer).List(m, &portsListServer{stream})
}

type Ports_ListServer interface {
	Send(*Port) error
	grpc.ServerStream
}

type portsListServer struct {
	grpc.ServerStream
}

func (x *portsListServer) Send(m *Port) error {
	return x.ServerStream.SendMsg(m)
}

var _Ports_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ports.v1.Ports",
	HandlerType: (*PortsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Ports_Get_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Update",
			Handler:       _Ports_Update_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "List",
			Handler:       _Ports_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "apis/ports/ports.proto",
}
