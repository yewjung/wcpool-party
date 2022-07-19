// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.2
// source: authorization.proto

package authorization

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuthorizationClient is the client API for Authorization service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthorizationClient interface {
	VerifyPartyID(ctx context.Context, in *Verification, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error)
}

type authorizationClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthorizationClient(cc grpc.ClientConnInterface) AuthorizationClient {
	return &authorizationClient{cc}
}

func (c *authorizationClient) VerifyPartyID(ctx context.Context, in *Verification, opts ...grpc.CallOption) (*wrapperspb.BoolValue, error) {
	out := new(wrapperspb.BoolValue)
	err := c.cc.Invoke(ctx, "/authorization.Authorization/VerifyPartyID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthorizationServer is the server API for Authorization service.
// All implementations must embed UnimplementedAuthorizationServer
// for forward compatibility
type AuthorizationServer interface {
	VerifyPartyID(context.Context, *Verification) (*wrapperspb.BoolValue, error)
	mustEmbedUnimplementedAuthorizationServer()
}

// UnimplementedAuthorizationServer must be embedded to have forward compatible implementations.
type UnimplementedAuthorizationServer struct {
}

func (UnimplementedAuthorizationServer) VerifyPartyID(context.Context, *Verification) (*wrapperspb.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyPartyID not implemented")
}
func (UnimplementedAuthorizationServer) mustEmbedUnimplementedAuthorizationServer() {}

// UnsafeAuthorizationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthorizationServer will
// result in compilation errors.
type UnsafeAuthorizationServer interface {
	mustEmbedUnimplementedAuthorizationServer()
}

func RegisterAuthorizationServer(s grpc.ServiceRegistrar, srv AuthorizationServer) {
	s.RegisterService(&Authorization_ServiceDesc, srv)
}

func _Authorization_VerifyPartyID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Verification)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthorizationServer).VerifyPartyID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/authorization.Authorization/VerifyPartyID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthorizationServer).VerifyPartyID(ctx, req.(*Verification))
	}
	return interceptor(ctx, in, info, handler)
}

// Authorization_ServiceDesc is the grpc.ServiceDesc for Authorization service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Authorization_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "authorization.Authorization",
	HandlerType: (*AuthorizationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VerifyPartyID",
			Handler:    _Authorization_VerifyPartyID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "authorization.proto",
}
