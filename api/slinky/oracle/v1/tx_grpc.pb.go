// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: slinky/oracle/v1/tx.proto

package oraclev1

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
	Msg_AddCurrencyPairs_FullMethodName    = "/slinky.oracle.v1.Msg/AddCurrencyPairs"
	Msg_RemoveCurrencyPairs_FullMethodName = "/slinky.oracle.v1.Msg/RemoveCurrencyPairs"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Msg is the message service for the x/oracle module.
type MsgClient interface {
	// AddCurrencyPairs will be used only by governance to update the set of
	// available CurrencyPairs. Given a set of CurrencyPair objects, update
	// the available currency pairs in the module .
	AddCurrencyPairs(ctx context.Context, in *MsgAddCurrencyPairs, opts ...grpc.CallOption) (*MsgAddCurrencyPairsResponse, error)
	// RemoveCurrencyPairs will be used explicitly by governance to remove the
	// given set of currency-pairs from the module's state. Thus these
	// CurrencyPairs will no longer have price-data available from this module.
	RemoveCurrencyPairs(ctx context.Context, in *MsgRemoveCurrencyPairs, opts ...grpc.CallOption) (*MsgRemoveCurrencyPairsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) AddCurrencyPairs(ctx context.Context, in *MsgAddCurrencyPairs, opts ...grpc.CallOption) (*MsgAddCurrencyPairsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgAddCurrencyPairsResponse)
	err := c.cc.Invoke(ctx, Msg_AddCurrencyPairs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RemoveCurrencyPairs(ctx context.Context, in *MsgRemoveCurrencyPairs, opts ...grpc.CallOption) (*MsgRemoveCurrencyPairsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgRemoveCurrencyPairsResponse)
	err := c.cc.Invoke(ctx, Msg_RemoveCurrencyPairs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility.
//
// Msg is the message service for the x/oracle module.
type MsgServer interface {
	// AddCurrencyPairs will be used only by governance to update the set of
	// available CurrencyPairs. Given a set of CurrencyPair objects, update
	// the available currency pairs in the module .
	AddCurrencyPairs(context.Context, *MsgAddCurrencyPairs) (*MsgAddCurrencyPairsResponse, error)
	// RemoveCurrencyPairs will be used explicitly by governance to remove the
	// given set of currency-pairs from the module's state. Thus these
	// CurrencyPairs will no longer have price-data available from this module.
	RemoveCurrencyPairs(context.Context, *MsgRemoveCurrencyPairs) (*MsgRemoveCurrencyPairsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMsgServer struct{}

func (UnimplementedMsgServer) AddCurrencyPairs(context.Context, *MsgAddCurrencyPairs) (*MsgAddCurrencyPairsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCurrencyPairs not implemented")
}
func (UnimplementedMsgServer) RemoveCurrencyPairs(context.Context, *MsgRemoveCurrencyPairs) (*MsgRemoveCurrencyPairsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveCurrencyPairs not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}
func (UnimplementedMsgServer) testEmbeddedByValue()             {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	// If the following call pancis, it indicates UnimplementedMsgServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_AddCurrencyPairs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddCurrencyPairs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddCurrencyPairs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_AddCurrencyPairs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddCurrencyPairs(ctx, req.(*MsgAddCurrencyPairs))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveCurrencyPairs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveCurrencyPairs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemoveCurrencyPairs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RemoveCurrencyPairs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemoveCurrencyPairs(ctx, req.(*MsgRemoveCurrencyPairs))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "slinky.oracle.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddCurrencyPairs",
			Handler:    _Msg_AddCurrencyPairs_Handler,
		},
		{
			MethodName: "RemoveCurrencyPairs",
			Handler:    _Msg_RemoveCurrencyPairs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "slinky/oracle/v1/tx.proto",
}