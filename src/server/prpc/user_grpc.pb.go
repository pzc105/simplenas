// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: user.proto

package prpc

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
	UserService_Register_FullMethodName         = "/prpc.UserService/Register"
	UserService_IsUsedEmail_FullMethodName      = "/prpc.UserService/IsUsedEmail"
	UserService_Login_FullMethodName            = "/prpc.UserService/Login"
	UserService_FastLogin_FullMethodName        = "/prpc.UserService/FastLogin"
	UserService_IsLogined_FullMethodName        = "/prpc.UserService/IsLogined"
	UserService_Download_FullMethodName         = "/prpc.UserService/Download"
	UserService_RemoveTorrent_FullMethodName    = "/prpc.UserService/RemoveTorrent"
	UserService_OnStatus_FullMethodName         = "/prpc.UserService/OnStatus"
	UserService_QueryBtVideos_FullMethodName    = "/prpc.UserService/QueryBtVideos"
	UserService_NewCategoryItem_FullMethodName  = "/prpc.UserService/NewCategoryItem"
	UserService_DelCategoryItem_FullMethodName  = "/prpc.UserService/DelCategoryItem"
	UserService_AddBtVideos_FullMethodName      = "/prpc.UserService/AddBtVideos"
	UserService_ShareItem_FullMethodName        = "/prpc.UserService/ShareItem"
	UserService_QuerySharedItems_FullMethodName = "/prpc.UserService/QuerySharedItems"
	UserService_DelSharedItem_FullMethodName    = "/prpc.UserService/DelSharedItem"
	UserService_QuerySubItems_FullMethodName    = "/prpc.UserService/QuerySubItems"
	UserService_QueryItemInfo_FullMethodName    = "/prpc.UserService/QueryItemInfo"
	UserService_RefreshSubtitle_FullMethodName  = "/prpc.UserService/RefreshSubtitle"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	Register(ctx context.Context, in *RegisterInfo, opts ...grpc.CallOption) (*RegisterRet, error)
	IsUsedEmail(ctx context.Context, in *EmailInfo, opts ...grpc.CallOption) (*IsUsedEmailRet, error)
	Login(ctx context.Context, in *LoginInfo, opts ...grpc.CallOption) (*LoginRet, error)
	FastLogin(ctx context.Context, in *LoginInfo, opts ...grpc.CallOption) (*LoginRet, error)
	IsLogined(ctx context.Context, in *LoginInfo, opts ...grpc.CallOption) (*LoginRet, error)
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadRespone, error)
	RemoveTorrent(ctx context.Context, in *RemoveTorrentReq, opts ...grpc.CallOption) (*RemoveTorrentRes, error)
	OnStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (UserService_OnStatusClient, error)
	QueryBtVideos(ctx context.Context, in *QueryBtVideosReq, opts ...grpc.CallOption) (*QueryBtVideosRes, error)
	NewCategoryItem(ctx context.Context, in *NewCategoryItemReq, opts ...grpc.CallOption) (*NewCategoryItemRes, error)
	DelCategoryItem(ctx context.Context, in *DelCategoryItemReq, opts ...grpc.CallOption) (*DelCategoryItemRes, error)
	AddBtVideos(ctx context.Context, in *AddBtVideosReq, opts ...grpc.CallOption) (*AddBtVideosRes, error)
	ShareItem(ctx context.Context, in *ShareItemReq, opts ...grpc.CallOption) (*ShareItemRes, error)
	QuerySharedItems(ctx context.Context, in *QuerySharedItemsReq, opts ...grpc.CallOption) (*QuerySharedItemsRes, error)
	DelSharedItem(ctx context.Context, in *DelSharedItemReq, opts ...grpc.CallOption) (*DelSharedItemRes, error)
	QuerySubItems(ctx context.Context, in *QuerySubItemsReq, opts ...grpc.CallOption) (*QuerySubItemsRes, error)
	QueryItemInfo(ctx context.Context, in *QueryItemInfoReq, opts ...grpc.CallOption) (*QueryItemInfoRes, error)
	RefreshSubtitle(ctx context.Context, in *RefreshSubtitleReq, opts ...grpc.CallOption) (*RefreshSubtitleRes, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Register(ctx context.Context, in *RegisterInfo, opts ...grpc.CallOption) (*RegisterRet, error) {
	out := new(RegisterRet)
	err := c.cc.Invoke(ctx, UserService_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) IsUsedEmail(ctx context.Context, in *EmailInfo, opts ...grpc.CallOption) (*IsUsedEmailRet, error) {
	out := new(IsUsedEmailRet)
	err := c.cc.Invoke(ctx, UserService_IsUsedEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Login(ctx context.Context, in *LoginInfo, opts ...grpc.CallOption) (*LoginRet, error) {
	out := new(LoginRet)
	err := c.cc.Invoke(ctx, UserService_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) FastLogin(ctx context.Context, in *LoginInfo, opts ...grpc.CallOption) (*LoginRet, error) {
	out := new(LoginRet)
	err := c.cc.Invoke(ctx, UserService_FastLogin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) IsLogined(ctx context.Context, in *LoginInfo, opts ...grpc.CallOption) (*LoginRet, error) {
	out := new(LoginRet)
	err := c.cc.Invoke(ctx, UserService_IsLogined_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadRespone, error) {
	out := new(DownloadRespone)
	err := c.cc.Invoke(ctx, UserService_Download_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RemoveTorrent(ctx context.Context, in *RemoveTorrentReq, opts ...grpc.CallOption) (*RemoveTorrentRes, error) {
	out := new(RemoveTorrentRes)
	err := c.cc.Invoke(ctx, UserService_RemoveTorrent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) OnStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (UserService_OnStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserService_ServiceDesc.Streams[0], UserService_OnStatus_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &userServiceOnStatusClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UserService_OnStatusClient interface {
	Recv() (*StatusRespone, error)
	grpc.ClientStream
}

type userServiceOnStatusClient struct {
	grpc.ClientStream
}

func (x *userServiceOnStatusClient) Recv() (*StatusRespone, error) {
	m := new(StatusRespone)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userServiceClient) QueryBtVideos(ctx context.Context, in *QueryBtVideosReq, opts ...grpc.CallOption) (*QueryBtVideosRes, error) {
	out := new(QueryBtVideosRes)
	err := c.cc.Invoke(ctx, UserService_QueryBtVideos_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) NewCategoryItem(ctx context.Context, in *NewCategoryItemReq, opts ...grpc.CallOption) (*NewCategoryItemRes, error) {
	out := new(NewCategoryItemRes)
	err := c.cc.Invoke(ctx, UserService_NewCategoryItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DelCategoryItem(ctx context.Context, in *DelCategoryItemReq, opts ...grpc.CallOption) (*DelCategoryItemRes, error) {
	out := new(DelCategoryItemRes)
	err := c.cc.Invoke(ctx, UserService_DelCategoryItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AddBtVideos(ctx context.Context, in *AddBtVideosReq, opts ...grpc.CallOption) (*AddBtVideosRes, error) {
	out := new(AddBtVideosRes)
	err := c.cc.Invoke(ctx, UserService_AddBtVideos_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ShareItem(ctx context.Context, in *ShareItemReq, opts ...grpc.CallOption) (*ShareItemRes, error) {
	out := new(ShareItemRes)
	err := c.cc.Invoke(ctx, UserService_ShareItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) QuerySharedItems(ctx context.Context, in *QuerySharedItemsReq, opts ...grpc.CallOption) (*QuerySharedItemsRes, error) {
	out := new(QuerySharedItemsRes)
	err := c.cc.Invoke(ctx, UserService_QuerySharedItems_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DelSharedItem(ctx context.Context, in *DelSharedItemReq, opts ...grpc.CallOption) (*DelSharedItemRes, error) {
	out := new(DelSharedItemRes)
	err := c.cc.Invoke(ctx, UserService_DelSharedItem_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) QuerySubItems(ctx context.Context, in *QuerySubItemsReq, opts ...grpc.CallOption) (*QuerySubItemsRes, error) {
	out := new(QuerySubItemsRes)
	err := c.cc.Invoke(ctx, UserService_QuerySubItems_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) QueryItemInfo(ctx context.Context, in *QueryItemInfoReq, opts ...grpc.CallOption) (*QueryItemInfoRes, error) {
	out := new(QueryItemInfoRes)
	err := c.cc.Invoke(ctx, UserService_QueryItemInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) RefreshSubtitle(ctx context.Context, in *RefreshSubtitleReq, opts ...grpc.CallOption) (*RefreshSubtitleRes, error) {
	out := new(RefreshSubtitleRes)
	err := c.cc.Invoke(ctx, UserService_RefreshSubtitle_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	Register(context.Context, *RegisterInfo) (*RegisterRet, error)
	IsUsedEmail(context.Context, *EmailInfo) (*IsUsedEmailRet, error)
	Login(context.Context, *LoginInfo) (*LoginRet, error)
	FastLogin(context.Context, *LoginInfo) (*LoginRet, error)
	IsLogined(context.Context, *LoginInfo) (*LoginRet, error)
	Download(context.Context, *DownloadRequest) (*DownloadRespone, error)
	RemoveTorrent(context.Context, *RemoveTorrentReq) (*RemoveTorrentRes, error)
	OnStatus(*StatusRequest, UserService_OnStatusServer) error
	QueryBtVideos(context.Context, *QueryBtVideosReq) (*QueryBtVideosRes, error)
	NewCategoryItem(context.Context, *NewCategoryItemReq) (*NewCategoryItemRes, error)
	DelCategoryItem(context.Context, *DelCategoryItemReq) (*DelCategoryItemRes, error)
	AddBtVideos(context.Context, *AddBtVideosReq) (*AddBtVideosRes, error)
	ShareItem(context.Context, *ShareItemReq) (*ShareItemRes, error)
	QuerySharedItems(context.Context, *QuerySharedItemsReq) (*QuerySharedItemsRes, error)
	DelSharedItem(context.Context, *DelSharedItemReq) (*DelSharedItemRes, error)
	QuerySubItems(context.Context, *QuerySubItemsReq) (*QuerySubItemsRes, error)
	QueryItemInfo(context.Context, *QueryItemInfoReq) (*QueryItemInfoRes, error)
	RefreshSubtitle(context.Context, *RefreshSubtitleReq) (*RefreshSubtitleRes, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) Register(context.Context, *RegisterInfo) (*RegisterRet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedUserServiceServer) IsUsedEmail(context.Context, *EmailInfo) (*IsUsedEmailRet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsUsedEmail not implemented")
}
func (UnimplementedUserServiceServer) Login(context.Context, *LoginInfo) (*LoginRet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedUserServiceServer) FastLogin(context.Context, *LoginInfo) (*LoginRet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FastLogin not implemented")
}
func (UnimplementedUserServiceServer) IsLogined(context.Context, *LoginInfo) (*LoginRet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsLogined not implemented")
}
func (UnimplementedUserServiceServer) Download(context.Context, *DownloadRequest) (*DownloadRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedUserServiceServer) RemoveTorrent(context.Context, *RemoveTorrentReq) (*RemoveTorrentRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveTorrent not implemented")
}
func (UnimplementedUserServiceServer) OnStatus(*StatusRequest, UserService_OnStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method OnStatus not implemented")
}
func (UnimplementedUserServiceServer) QueryBtVideos(context.Context, *QueryBtVideosReq) (*QueryBtVideosRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryBtVideos not implemented")
}
func (UnimplementedUserServiceServer) NewCategoryItem(context.Context, *NewCategoryItemReq) (*NewCategoryItemRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewCategoryItem not implemented")
}
func (UnimplementedUserServiceServer) DelCategoryItem(context.Context, *DelCategoryItemReq) (*DelCategoryItemRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelCategoryItem not implemented")
}
func (UnimplementedUserServiceServer) AddBtVideos(context.Context, *AddBtVideosReq) (*AddBtVideosRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddBtVideos not implemented")
}
func (UnimplementedUserServiceServer) ShareItem(context.Context, *ShareItemReq) (*ShareItemRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShareItem not implemented")
}
func (UnimplementedUserServiceServer) QuerySharedItems(context.Context, *QuerySharedItemsReq) (*QuerySharedItemsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuerySharedItems not implemented")
}
func (UnimplementedUserServiceServer) DelSharedItem(context.Context, *DelSharedItemReq) (*DelSharedItemRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelSharedItem not implemented")
}
func (UnimplementedUserServiceServer) QuerySubItems(context.Context, *QuerySubItemsReq) (*QuerySubItemsRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuerySubItems not implemented")
}
func (UnimplementedUserServiceServer) QueryItemInfo(context.Context, *QueryItemInfoReq) (*QueryItemInfoRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryItemInfo not implemented")
}
func (UnimplementedUserServiceServer) RefreshSubtitle(context.Context, *RefreshSubtitleReq) (*RefreshSubtitleRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshSubtitle not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Register(ctx, req.(*RegisterInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_IsUsedEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).IsUsedEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_IsUsedEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).IsUsedEmail(ctx, req.(*EmailInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Login(ctx, req.(*LoginInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_FastLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).FastLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_FastLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).FastLogin(ctx, req.(*LoginInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_IsLogined_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).IsLogined(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_IsLogined_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).IsLogined(ctx, req.(*LoginInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_Download_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Download(ctx, req.(*DownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RemoveTorrent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveTorrentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RemoveTorrent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_RemoveTorrent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RemoveTorrent(ctx, req.(*RemoveTorrentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_OnStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StatusRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServiceServer).OnStatus(m, &userServiceOnStatusServer{stream})
}

type UserService_OnStatusServer interface {
	Send(*StatusRespone) error
	grpc.ServerStream
}

type userServiceOnStatusServer struct {
	grpc.ServerStream
}

func (x *userServiceOnStatusServer) Send(m *StatusRespone) error {
	return x.ServerStream.SendMsg(m)
}

func _UserService_QueryBtVideos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBtVideosReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).QueryBtVideos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_QueryBtVideos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).QueryBtVideos(ctx, req.(*QueryBtVideosReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_NewCategoryItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewCategoryItemReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).NewCategoryItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_NewCategoryItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).NewCategoryItem(ctx, req.(*NewCategoryItemReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DelCategoryItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelCategoryItemReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DelCategoryItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_DelCategoryItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DelCategoryItem(ctx, req.(*DelCategoryItemReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AddBtVideos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddBtVideosReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AddBtVideos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_AddBtVideos_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AddBtVideos(ctx, req.(*AddBtVideosReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ShareItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareItemReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ShareItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_ShareItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ShareItem(ctx, req.(*ShareItemReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_QuerySharedItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySharedItemsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).QuerySharedItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_QuerySharedItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).QuerySharedItems(ctx, req.(*QuerySharedItemsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DelSharedItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelSharedItemReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DelSharedItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_DelSharedItem_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DelSharedItem(ctx, req.(*DelSharedItemReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_QuerySubItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuerySubItemsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).QuerySubItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_QuerySubItems_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).QuerySubItems(ctx, req.(*QuerySubItemsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_QueryItemInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryItemInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).QueryItemInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_QueryItemInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).QueryItemInfo(ctx, req.(*QueryItemInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_RefreshSubtitle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshSubtitleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).RefreshSubtitle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_RefreshSubtitle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).RefreshSubtitle(ctx, req.(*RefreshSubtitleReq))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "prpc.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _UserService_Register_Handler,
		},
		{
			MethodName: "IsUsedEmail",
			Handler:    _UserService_IsUsedEmail_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _UserService_Login_Handler,
		},
		{
			MethodName: "FastLogin",
			Handler:    _UserService_FastLogin_Handler,
		},
		{
			MethodName: "IsLogined",
			Handler:    _UserService_IsLogined_Handler,
		},
		{
			MethodName: "Download",
			Handler:    _UserService_Download_Handler,
		},
		{
			MethodName: "RemoveTorrent",
			Handler:    _UserService_RemoveTorrent_Handler,
		},
		{
			MethodName: "QueryBtVideos",
			Handler:    _UserService_QueryBtVideos_Handler,
		},
		{
			MethodName: "NewCategoryItem",
			Handler:    _UserService_NewCategoryItem_Handler,
		},
		{
			MethodName: "DelCategoryItem",
			Handler:    _UserService_DelCategoryItem_Handler,
		},
		{
			MethodName: "AddBtVideos",
			Handler:    _UserService_AddBtVideos_Handler,
		},
		{
			MethodName: "ShareItem",
			Handler:    _UserService_ShareItem_Handler,
		},
		{
			MethodName: "QuerySharedItems",
			Handler:    _UserService_QuerySharedItems_Handler,
		},
		{
			MethodName: "DelSharedItem",
			Handler:    _UserService_DelSharedItem_Handler,
		},
		{
			MethodName: "QuerySubItems",
			Handler:    _UserService_QuerySubItems_Handler,
		},
		{
			MethodName: "QueryItemInfo",
			Handler:    _UserService_QueryItemInfo_Handler,
		},
		{
			MethodName: "RefreshSubtitle",
			Handler:    _UserService_RefreshSubtitle_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "OnStatus",
			Handler:       _UserService_OnStatus_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "user.proto",
}
