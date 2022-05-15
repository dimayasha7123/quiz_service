// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: api/api.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// QuizServiceClient is the client API for QuizService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QuizServiceClient interface {
	GetQuizList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*QuizList, error)
	CreateQuizParty(ctx context.Context, in *QuizParty, opts ...grpc.CallOption) (*QuizPartyID, error)
	GetNextQuestion(ctx context.Context, in *QuizPartyID, opts ...grpc.CallOption) (*QuestionOrNil, error)
	SendAnswer(ctx context.Context, in *Answer, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetSingleTop(ctx context.Context, in *QuizPartyID, opts ...grpc.CallOption) (*SingleTop, error)
	GetGlobalQuizTop(ctx context.Context, in *QuizID, opts ...grpc.CallOption) (*GlobalTop, error)
}

type quizServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQuizServiceClient(cc grpc.ClientConnInterface) QuizServiceClient {
	return &quizServiceClient{cc}
}

func (c *quizServiceClient) GetQuizList(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*QuizList, error) {
	out := new(QuizList)
	err := c.cc.Invoke(ctx, "/api.QuizService/GetQuizList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizServiceClient) CreateQuizParty(ctx context.Context, in *QuizParty, opts ...grpc.CallOption) (*QuizPartyID, error) {
	out := new(QuizPartyID)
	err := c.cc.Invoke(ctx, "/api.QuizService/CreateQuizParty", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizServiceClient) GetNextQuestion(ctx context.Context, in *QuizPartyID, opts ...grpc.CallOption) (*QuestionOrNil, error) {
	out := new(QuestionOrNil)
	err := c.cc.Invoke(ctx, "/api.QuizService/GetNextQuestion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizServiceClient) SendAnswer(ctx context.Context, in *Answer, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/api.QuizService/SendAnswer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizServiceClient) GetSingleTop(ctx context.Context, in *QuizPartyID, opts ...grpc.CallOption) (*SingleTop, error) {
	out := new(SingleTop)
	err := c.cc.Invoke(ctx, "/api.QuizService/GetSingleTop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizServiceClient) GetGlobalQuizTop(ctx context.Context, in *QuizID, opts ...grpc.CallOption) (*GlobalTop, error) {
	out := new(GlobalTop)
	err := c.cc.Invoke(ctx, "/api.QuizService/GetGlobalQuizTop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QuizServiceServer is the server API for QuizService service.
// All implementations must embed UnimplementedQuizServiceServer
// for forward compatibility
type QuizServiceServer interface {
	GetQuizList(context.Context, *emptypb.Empty) (*QuizList, error)
	CreateQuizParty(context.Context, *QuizParty) (*QuizPartyID, error)
	GetNextQuestion(context.Context, *QuizPartyID) (*QuestionOrNil, error)
	SendAnswer(context.Context, *Answer) (*emptypb.Empty, error)
	GetSingleTop(context.Context, *QuizPartyID) (*SingleTop, error)
	GetGlobalQuizTop(context.Context, *QuizID) (*GlobalTop, error)
	mustEmbedUnimplementedQuizServiceServer()
}

// UnimplementedQuizServiceServer must be embedded to have forward compatible implementations.
type UnimplementedQuizServiceServer struct {
}

func (UnimplementedQuizServiceServer) GetQuizList(context.Context, *emptypb.Empty) (*QuizList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuizList not implemented")
}
func (UnimplementedQuizServiceServer) CreateQuizParty(context.Context, *QuizParty) (*QuizPartyID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateQuizParty not implemented")
}
func (UnimplementedQuizServiceServer) GetNextQuestion(context.Context, *QuizPartyID) (*QuestionOrNil, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNextQuestion not implemented")
}
func (UnimplementedQuizServiceServer) SendAnswer(context.Context, *Answer) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendAnswer not implemented")
}
func (UnimplementedQuizServiceServer) GetSingleTop(context.Context, *QuizPartyID) (*SingleTop, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSingleTop not implemented")
}
func (UnimplementedQuizServiceServer) GetGlobalQuizTop(context.Context, *QuizID) (*GlobalTop, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGlobalQuizTop not implemented")
}
func (UnimplementedQuizServiceServer) mustEmbedUnimplementedQuizServiceServer() {}

// UnsafeQuizServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QuizServiceServer will
// result in compilation errors.
type UnsafeQuizServiceServer interface {
	mustEmbedUnimplementedQuizServiceServer()
}

func RegisterQuizServiceServer(s grpc.ServiceRegistrar, srv QuizServiceServer) {
	s.RegisterService(&QuizService_ServiceDesc, srv)
}

func _QuizService_GetQuizList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuizServiceServer).GetQuizList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.QuizService/GetQuizList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuizServiceServer).GetQuizList(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuizService_CreateQuizParty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuizParty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuizServiceServer).CreateQuizParty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.QuizService/CreateQuizParty",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuizServiceServer).CreateQuizParty(ctx, req.(*QuizParty))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuizService_GetNextQuestion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuizPartyID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuizServiceServer).GetNextQuestion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.QuizService/GetNextQuestion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuizServiceServer).GetNextQuestion(ctx, req.(*QuizPartyID))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuizService_SendAnswer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Answer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuizServiceServer).SendAnswer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.QuizService/SendAnswer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuizServiceServer).SendAnswer(ctx, req.(*Answer))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuizService_GetSingleTop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuizPartyID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuizServiceServer).GetSingleTop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.QuizService/GetSingleTop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuizServiceServer).GetSingleTop(ctx, req.(*QuizPartyID))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuizService_GetGlobalQuizTop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuizID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuizServiceServer).GetGlobalQuizTop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.QuizService/GetGlobalQuizTop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuizServiceServer).GetGlobalQuizTop(ctx, req.(*QuizID))
	}
	return interceptor(ctx, in, info, handler)
}

// QuizService_ServiceDesc is the grpc.ServiceDesc for QuizService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QuizService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.QuizService",
	HandlerType: (*QuizServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetQuizList",
			Handler:    _QuizService_GetQuizList_Handler,
		},
		{
			MethodName: "CreateQuizParty",
			Handler:    _QuizService_CreateQuizParty_Handler,
		},
		{
			MethodName: "GetNextQuestion",
			Handler:    _QuizService_GetNextQuestion_Handler,
		},
		{
			MethodName: "SendAnswer",
			Handler:    _QuizService_SendAnswer_Handler,
		},
		{
			MethodName: "GetSingleTop",
			Handler:    _QuizService_GetSingleTop_Handler,
		},
		{
			MethodName: "GetGlobalQuizTop",
			Handler:    _QuizService_GetGlobalQuizTop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/api.proto",
}
