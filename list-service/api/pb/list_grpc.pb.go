// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.15.8
// source: list.proto

package list

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

// ListServiceClient is the client API for ListService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ListServiceClient interface {
	CreateList(ctx context.Context, in *CreateListRequest, opts ...grpc.CallOption) (*CreateListResponse, error)
	GetListByID(ctx context.Context, in *GetListByIDRequest, opts ...grpc.CallOption) (*GetListByIDResponse, error)
	GetListsByBoard(ctx context.Context, in *GetListsByBoardRequest, opts ...grpc.CallOption) (*GetListsByBoardResponse, error)
	UpdateListName(ctx context.Context, in *UpdateListNameRequest, opts ...grpc.CallOption) (*UpdateListNameResponse, error)
	MoveListPosition(ctx context.Context, in *MoveListPositionRequest, opts ...grpc.CallOption) (*MoveListPositionResponse, error)
	ArchiveList(ctx context.Context, in *ArchiveListRequest, opts ...grpc.CallOption) (*ArchiveListResponse, error)
	RestoreList(ctx context.Context, in *RestoreListRequest, opts ...grpc.CallOption) (*RestoreListResponse, error)
	DeleteList(ctx context.Context, in *DeleteListRequest, opts ...grpc.CallOption) (*DeleteListResponse, error)
}

type listServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewListServiceClient(cc grpc.ClientConnInterface) ListServiceClient {
	return &listServiceClient{cc}
}

func (c *listServiceClient) CreateList(ctx context.Context, in *CreateListRequest, opts ...grpc.CallOption) (*CreateListResponse, error) {
	out := new(CreateListResponse)
	err := c.cc.Invoke(ctx, "/listpb.ListService/CreateList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listServiceClient) GetListByID(ctx context.Context, in *GetListByIDRequest, opts ...grpc.CallOption) (*GetListByIDResponse, error) {
	out := new(GetListByIDResponse)
	err := c.cc.Invoke(ctx, "/listpb.ListService/GetListByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listServiceClient) GetListsByBoard(ctx context.Context, in *GetListsByBoardRequest, opts ...grpc.CallOption) (*GetListsByBoardResponse, error) {
	out := new(GetListsByBoardResponse)
	err := c.cc.Invoke(ctx, "/listpb.ListService/GetListsByBoard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listServiceClient) UpdateListName(ctx context.Context, in *UpdateListNameRequest, opts ...grpc.CallOption) (*UpdateListNameResponse, error) {
	out := new(UpdateListNameResponse)
	err := c.cc.Invoke(ctx, "/listpb.ListService/UpdateListName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listServiceClient) MoveListPosition(ctx context.Context, in *MoveListPositionRequest, opts ...grpc.CallOption) (*MoveListPositionResponse, error) {
	out := new(MoveListPositionResponse)
	err := c.cc.Invoke(ctx, "/listpb.ListService/MoveListPosition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listServiceClient) ArchiveList(ctx context.Context, in *ArchiveListRequest, opts ...grpc.CallOption) (*ArchiveListResponse, error) {
	out := new(ArchiveListResponse)
	err := c.cc.Invoke(ctx, "/listpb.ListService/ArchiveList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listServiceClient) RestoreList(ctx context.Context, in *RestoreListRequest, opts ...grpc.CallOption) (*RestoreListResponse, error) {
	out := new(RestoreListResponse)
	err := c.cc.Invoke(ctx, "/listpb.ListService/RestoreList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *listServiceClient) DeleteList(ctx context.Context, in *DeleteListRequest, opts ...grpc.CallOption) (*DeleteListResponse, error) {
	out := new(DeleteListResponse)
	err := c.cc.Invoke(ctx, "/listpb.ListService/DeleteList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListServiceServer is the server API for ListService service.
// All implementations must embed UnimplementedListServiceServer
// for forward compatibility
type ListServiceServer interface {
	CreateList(context.Context, *CreateListRequest) (*CreateListResponse, error)
	GetListByID(context.Context, *GetListByIDRequest) (*GetListByIDResponse, error)
	GetListsByBoard(context.Context, *GetListsByBoardRequest) (*GetListsByBoardResponse, error)
	UpdateListName(context.Context, *UpdateListNameRequest) (*UpdateListNameResponse, error)
	MoveListPosition(context.Context, *MoveListPositionRequest) (*MoveListPositionResponse, error)
	ArchiveList(context.Context, *ArchiveListRequest) (*ArchiveListResponse, error)
	RestoreList(context.Context, *RestoreListRequest) (*RestoreListResponse, error)
	DeleteList(context.Context, *DeleteListRequest) (*DeleteListResponse, error)
	mustEmbedUnimplementedListServiceServer()
}

// UnimplementedListServiceServer must be embedded to have forward compatible implementations.
type UnimplementedListServiceServer struct {
}

func (UnimplementedListServiceServer) CreateList(context.Context, *CreateListRequest) (*CreateListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateList not implemented")
}
func (UnimplementedListServiceServer) GetListByID(context.Context, *GetListByIDRequest) (*GetListByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListByID not implemented")
}
func (UnimplementedListServiceServer) GetListsByBoard(context.Context, *GetListsByBoardRequest) (*GetListsByBoardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListsByBoard not implemented")
}
func (UnimplementedListServiceServer) UpdateListName(context.Context, *UpdateListNameRequest) (*UpdateListNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateListName not implemented")
}
func (UnimplementedListServiceServer) MoveListPosition(context.Context, *MoveListPositionRequest) (*MoveListPositionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoveListPosition not implemented")
}
func (UnimplementedListServiceServer) ArchiveList(context.Context, *ArchiveListRequest) (*ArchiveListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ArchiveList not implemented")
}
func (UnimplementedListServiceServer) RestoreList(context.Context, *RestoreListRequest) (*RestoreListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestoreList not implemented")
}
func (UnimplementedListServiceServer) DeleteList(context.Context, *DeleteListRequest) (*DeleteListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteList not implemented")
}
func (UnimplementedListServiceServer) mustEmbedUnimplementedListServiceServer() {}

// UnsafeListServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ListServiceServer will
// result in compilation errors.
type UnsafeListServiceServer interface {
	mustEmbedUnimplementedListServiceServer()
}

func RegisterListServiceServer(s grpc.ServiceRegistrar, srv ListServiceServer) {
	s.RegisterService(&ListService_ServiceDesc, srv)
}

func _ListService_CreateList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).CreateList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/listpb.ListService/CreateList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).CreateList(ctx, req.(*CreateListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListService_GetListByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).GetListByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/listpb.ListService/GetListByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).GetListByID(ctx, req.(*GetListByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListService_GetListsByBoard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListsByBoardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).GetListsByBoard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/listpb.ListService/GetListsByBoard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).GetListsByBoard(ctx, req.(*GetListsByBoardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListService_UpdateListName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateListNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).UpdateListName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/listpb.ListService/UpdateListName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).UpdateListName(ctx, req.(*UpdateListNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListService_MoveListPosition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MoveListPositionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).MoveListPosition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/listpb.ListService/MoveListPosition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).MoveListPosition(ctx, req.(*MoveListPositionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListService_ArchiveList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ArchiveListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).ArchiveList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/listpb.ListService/ArchiveList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).ArchiveList(ctx, req.(*ArchiveListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListService_RestoreList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RestoreListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).RestoreList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/listpb.ListService/RestoreList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).RestoreList(ctx, req.(*RestoreListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ListService_DeleteList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ListServiceServer).DeleteList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/listpb.ListService/DeleteList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ListServiceServer).DeleteList(ctx, req.(*DeleteListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ListService_ServiceDesc is the grpc.ServiceDesc for ListService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ListService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "listpb.ListService",
	HandlerType: (*ListServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateList",
			Handler:    _ListService_CreateList_Handler,
		},
		{
			MethodName: "GetListByID",
			Handler:    _ListService_GetListByID_Handler,
		},
		{
			MethodName: "GetListsByBoard",
			Handler:    _ListService_GetListsByBoard_Handler,
		},
		{
			MethodName: "UpdateListName",
			Handler:    _ListService_UpdateListName_Handler,
		},
		{
			MethodName: "MoveListPosition",
			Handler:    _ListService_MoveListPosition_Handler,
		},
		{
			MethodName: "ArchiveList",
			Handler:    _ListService_ArchiveList_Handler,
		},
		{
			MethodName: "RestoreList",
			Handler:    _ListService_RestoreList_Handler,
		},
		{
			MethodName: "DeleteList",
			Handler:    _ListService_DeleteList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "list.proto",
}
