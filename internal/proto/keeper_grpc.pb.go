// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: internal/proto/keeper.proto

package proto

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
	Keeper_Ping_FullMethodName                 = "/proto.Keeper/Ping"
	Keeper_Registration_FullMethodName         = "/proto.Keeper/Registration"
	Keeper_Login_FullMethodName                = "/proto.Keeper/Login"
	Keeper_EntityCodes_FullMethodName          = "/proto.Keeper/EntityCodes"
	Keeper_Fields_FullMethodName               = "/proto.Keeper/Fields"
	Keeper_AddEntity_FullMethodName            = "/proto.Keeper/AddEntity"
	Keeper_UploadBinary_FullMethodName         = "/proto.Keeper/UploadBinary"
	Keeper_UploadCryptoBinary_FullMethodName   = "/proto.Keeper/UploadCryptoBinary"
	Keeper_Entity_FullMethodName               = "/proto.Keeper/Entity"
	Keeper_DownloadBinary_FullMethodName       = "/proto.Keeper/DownloadBinary"
	Keeper_DownloadCryptoBinary_FullMethodName = "/proto.Keeper/DownloadCryptoBinary"
	Keeper_EntityList_FullMethodName           = "/proto.Keeper/EntityList"
)

// KeeperClient is the client API for Keeper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeeperClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Registration(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	EntityCodes(ctx context.Context, in *EntityCodesRequest, opts ...grpc.CallOption) (*EntityCodesResponse, error)
	Fields(ctx context.Context, in *FieldsRequest, opts ...grpc.CallOption) (*FieldsResponse, error)
	AddEntity(ctx context.Context, in *AddEntityRequest, opts ...grpc.CallOption) (*AddEntityResponse, error)
	UploadBinary(ctx context.Context, opts ...grpc.CallOption) (Keeper_UploadBinaryClient, error)
	UploadCryptoBinary(ctx context.Context, opts ...grpc.CallOption) (Keeper_UploadCryptoBinaryClient, error)
	Entity(ctx context.Context, in *EntityRequest, opts ...grpc.CallOption) (*EntityResponse, error)
	DownloadBinary(ctx context.Context, in *DownloadBinRequest, opts ...grpc.CallOption) (Keeper_DownloadBinaryClient, error)
	DownloadCryptoBinary(ctx context.Context, in *DownloadBinRequest, opts ...grpc.CallOption) (Keeper_DownloadCryptoBinaryClient, error)
	EntityList(ctx context.Context, in *EntityListRequest, opts ...grpc.CallOption) (*EntityListResponse, error)
}

type keeperClient struct {
	cc grpc.ClientConnInterface
}

func NewKeeperClient(cc grpc.ClientConnInterface) KeeperClient {
	return &keeperClient{cc}
}

func (c *keeperClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, Keeper_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) Registration(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, Keeper_Registration_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, Keeper_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) EntityCodes(ctx context.Context, in *EntityCodesRequest, opts ...grpc.CallOption) (*EntityCodesResponse, error) {
	out := new(EntityCodesResponse)
	err := c.cc.Invoke(ctx, Keeper_EntityCodes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) Fields(ctx context.Context, in *FieldsRequest, opts ...grpc.CallOption) (*FieldsResponse, error) {
	out := new(FieldsResponse)
	err := c.cc.Invoke(ctx, Keeper_Fields_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) AddEntity(ctx context.Context, in *AddEntityRequest, opts ...grpc.CallOption) (*AddEntityResponse, error) {
	out := new(AddEntityResponse)
	err := c.cc.Invoke(ctx, Keeper_AddEntity_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) UploadBinary(ctx context.Context, opts ...grpc.CallOption) (Keeper_UploadBinaryClient, error) {
	stream, err := c.cc.NewStream(ctx, &Keeper_ServiceDesc.Streams[0], Keeper_UploadBinary_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &keeperUploadBinaryClient{stream}
	return x, nil
}

type Keeper_UploadBinaryClient interface {
	Send(*UploadBinRequest) error
	CloseAndRecv() (*UploadBinResponse, error)
	grpc.ClientStream
}

type keeperUploadBinaryClient struct {
	grpc.ClientStream
}

func (x *keeperUploadBinaryClient) Send(m *UploadBinRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *keeperUploadBinaryClient) CloseAndRecv() (*UploadBinResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadBinResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *keeperClient) UploadCryptoBinary(ctx context.Context, opts ...grpc.CallOption) (Keeper_UploadCryptoBinaryClient, error) {
	stream, err := c.cc.NewStream(ctx, &Keeper_ServiceDesc.Streams[1], Keeper_UploadCryptoBinary_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &keeperUploadCryptoBinaryClient{stream}
	return x, nil
}

type Keeper_UploadCryptoBinaryClient interface {
	Send(*UploadBinRequest) error
	CloseAndRecv() (*UploadBinResponse, error)
	grpc.ClientStream
}

type keeperUploadCryptoBinaryClient struct {
	grpc.ClientStream
}

func (x *keeperUploadCryptoBinaryClient) Send(m *UploadBinRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *keeperUploadCryptoBinaryClient) CloseAndRecv() (*UploadBinResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadBinResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *keeperClient) Entity(ctx context.Context, in *EntityRequest, opts ...grpc.CallOption) (*EntityResponse, error) {
	out := new(EntityResponse)
	err := c.cc.Invoke(ctx, Keeper_Entity_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperClient) DownloadBinary(ctx context.Context, in *DownloadBinRequest, opts ...grpc.CallOption) (Keeper_DownloadBinaryClient, error) {
	stream, err := c.cc.NewStream(ctx, &Keeper_ServiceDesc.Streams[2], Keeper_DownloadBinary_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &keeperDownloadBinaryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Keeper_DownloadBinaryClient interface {
	Recv() (*DownloadBinResponse, error)
	grpc.ClientStream
}

type keeperDownloadBinaryClient struct {
	grpc.ClientStream
}

func (x *keeperDownloadBinaryClient) Recv() (*DownloadBinResponse, error) {
	m := new(DownloadBinResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *keeperClient) DownloadCryptoBinary(ctx context.Context, in *DownloadBinRequest, opts ...grpc.CallOption) (Keeper_DownloadCryptoBinaryClient, error) {
	stream, err := c.cc.NewStream(ctx, &Keeper_ServiceDesc.Streams[3], Keeper_DownloadCryptoBinary_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &keeperDownloadCryptoBinaryClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Keeper_DownloadCryptoBinaryClient interface {
	Recv() (*DownloadBinResponse, error)
	grpc.ClientStream
}

type keeperDownloadCryptoBinaryClient struct {
	grpc.ClientStream
}

func (x *keeperDownloadCryptoBinaryClient) Recv() (*DownloadBinResponse, error) {
	m := new(DownloadBinResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *keeperClient) EntityList(ctx context.Context, in *EntityListRequest, opts ...grpc.CallOption) (*EntityListResponse, error) {
	out := new(EntityListResponse)
	err := c.cc.Invoke(ctx, Keeper_EntityList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeeperServer is the server API for Keeper service.
// All implementations must embed UnimplementedKeeperServer
// for forward compatibility
type KeeperServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Registration(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	EntityCodes(context.Context, *EntityCodesRequest) (*EntityCodesResponse, error)
	Fields(context.Context, *FieldsRequest) (*FieldsResponse, error)
	AddEntity(context.Context, *AddEntityRequest) (*AddEntityResponse, error)
	UploadBinary(Keeper_UploadBinaryServer) error
	UploadCryptoBinary(Keeper_UploadCryptoBinaryServer) error
	Entity(context.Context, *EntityRequest) (*EntityResponse, error)
	DownloadBinary(*DownloadBinRequest, Keeper_DownloadBinaryServer) error
	DownloadCryptoBinary(*DownloadBinRequest, Keeper_DownloadCryptoBinaryServer) error
	EntityList(context.Context, *EntityListRequest) (*EntityListResponse, error)
	mustEmbedUnimplementedKeeperServer()
}

// UnimplementedKeeperServer must be embedded to have forward compatible implementations.
type UnimplementedKeeperServer struct {
}

func (UnimplementedKeeperServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedKeeperServer) Registration(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Registration not implemented")
}
func (UnimplementedKeeperServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedKeeperServer) EntityCodes(context.Context, *EntityCodesRequest) (*EntityCodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EntityCodes not implemented")
}
func (UnimplementedKeeperServer) Fields(context.Context, *FieldsRequest) (*FieldsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fields not implemented")
}
func (UnimplementedKeeperServer) AddEntity(context.Context, *AddEntityRequest) (*AddEntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddEntity not implemented")
}
func (UnimplementedKeeperServer) UploadBinary(Keeper_UploadBinaryServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadBinary not implemented")
}
func (UnimplementedKeeperServer) UploadCryptoBinary(Keeper_UploadCryptoBinaryServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadCryptoBinary not implemented")
}
func (UnimplementedKeeperServer) Entity(context.Context, *EntityRequest) (*EntityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Entity not implemented")
}
func (UnimplementedKeeperServer) DownloadBinary(*DownloadBinRequest, Keeper_DownloadBinaryServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadBinary not implemented")
}
func (UnimplementedKeeperServer) DownloadCryptoBinary(*DownloadBinRequest, Keeper_DownloadCryptoBinaryServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadCryptoBinary not implemented")
}
func (UnimplementedKeeperServer) EntityList(context.Context, *EntityListRequest) (*EntityListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EntityList not implemented")
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

func _Keeper_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_Registration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).Registration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_Registration_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).Registration(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_EntityCodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntityCodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).EntityCodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_EntityCodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).EntityCodes(ctx, req.(*EntityCodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_Fields_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FieldsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).Fields(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_Fields_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).Fields(ctx, req.(*FieldsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_AddEntity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddEntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).AddEntity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_AddEntity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).AddEntity(ctx, req.(*AddEntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_UploadBinary_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(KeeperServer).UploadBinary(&keeperUploadBinaryServer{stream})
}

type Keeper_UploadBinaryServer interface {
	SendAndClose(*UploadBinResponse) error
	Recv() (*UploadBinRequest, error)
	grpc.ServerStream
}

type keeperUploadBinaryServer struct {
	grpc.ServerStream
}

func (x *keeperUploadBinaryServer) SendAndClose(m *UploadBinResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *keeperUploadBinaryServer) Recv() (*UploadBinRequest, error) {
	m := new(UploadBinRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Keeper_UploadCryptoBinary_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(KeeperServer).UploadCryptoBinary(&keeperUploadCryptoBinaryServer{stream})
}

type Keeper_UploadCryptoBinaryServer interface {
	SendAndClose(*UploadBinResponse) error
	Recv() (*UploadBinRequest, error)
	grpc.ServerStream
}

type keeperUploadCryptoBinaryServer struct {
	grpc.ServerStream
}

func (x *keeperUploadCryptoBinaryServer) SendAndClose(m *UploadBinResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *keeperUploadCryptoBinaryServer) Recv() (*UploadBinRequest, error) {
	m := new(UploadBinRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Keeper_Entity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).Entity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_Entity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).Entity(ctx, req.(*EntityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Keeper_DownloadBinary_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadBinRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KeeperServer).DownloadBinary(m, &keeperDownloadBinaryServer{stream})
}

type Keeper_DownloadBinaryServer interface {
	Send(*DownloadBinResponse) error
	grpc.ServerStream
}

type keeperDownloadBinaryServer struct {
	grpc.ServerStream
}

func (x *keeperDownloadBinaryServer) Send(m *DownloadBinResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Keeper_DownloadCryptoBinary_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadBinRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KeeperServer).DownloadCryptoBinary(m, &keeperDownloadCryptoBinaryServer{stream})
}

type Keeper_DownloadCryptoBinaryServer interface {
	Send(*DownloadBinResponse) error
	grpc.ServerStream
}

type keeperDownloadCryptoBinaryServer struct {
	grpc.ServerStream
}

func (x *keeperDownloadCryptoBinaryServer) Send(m *DownloadBinResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Keeper_EntityList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EntityListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServer).EntityList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Keeper_EntityList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServer).EntityList(ctx, req.(*EntityListRequest))
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
			MethodName: "Ping",
			Handler:    _Keeper_Ping_Handler,
		},
		{
			MethodName: "Registration",
			Handler:    _Keeper_Registration_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Keeper_Login_Handler,
		},
		{
			MethodName: "EntityCodes",
			Handler:    _Keeper_EntityCodes_Handler,
		},
		{
			MethodName: "Fields",
			Handler:    _Keeper_Fields_Handler,
		},
		{
			MethodName: "AddEntity",
			Handler:    _Keeper_AddEntity_Handler,
		},
		{
			MethodName: "Entity",
			Handler:    _Keeper_Entity_Handler,
		},
		{
			MethodName: "EntityList",
			Handler:    _Keeper_EntityList_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadBinary",
			Handler:       _Keeper_UploadBinary_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "UploadCryptoBinary",
			Handler:       _Keeper_UploadCryptoBinary_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadBinary",
			Handler:       _Keeper_DownloadBinary_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DownloadCryptoBinary",
			Handler:       _Keeper_DownloadCryptoBinary_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/proto/keeper.proto",
}
