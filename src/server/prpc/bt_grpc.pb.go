// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: bt.proto

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
	BtService_InitedSession_FullMethodName    = "/prpc.BtService/InitedSession"
	BtService_InitSession_FullMethodName      = "/prpc.BtService/InitSession"
	BtService_Parse_FullMethodName            = "/prpc.BtService/Parse"
	BtService_Download_FullMethodName         = "/prpc.BtService/Download"
	BtService_RemoveTorrent_FullMethodName    = "/prpc.BtService/RemoveTorrent"
	BtService_GetMagnetUri_FullMethodName     = "/prpc.BtService/GetMagnetUri"
	BtService_GetResumeData_FullMethodName    = "/prpc.BtService/GetResumeData"
	BtService_GetTorrentInfo_FullMethodName   = "/prpc.BtService/GetTorrentInfo"
	BtService_GetBtStatus_FullMethodName      = "/prpc.BtService/GetBtStatus"
	BtService_GetSessionParams_FullMethodName = "/prpc.BtService/GetSessionParams"
	BtService_GetPeerInfo_FullMethodName      = "/prpc.BtService/GetPeerInfo"
	BtService_OnBtStatus_FullMethodName       = "/prpc.BtService/OnBtStatus"
	BtService_OnFileCompleted_FullMethodName  = "/prpc.BtService/OnFileCompleted"
)

// BtServiceClient is the client API for BtService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BtServiceClient interface {
	InitedSession(ctx context.Context, in *InitedSessionReq, opts ...grpc.CallOption) (*InitedSessionRsp, error)
	InitSession(ctx context.Context, in *InitSessionReq, opts ...grpc.CallOption) (*InitSessionRsp, error)
	Parse(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadRespone, error)
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadRespone, error)
	RemoveTorrent(ctx context.Context, in *RemoveTorrentReq, opts ...grpc.CallOption) (*RemoveTorrentRes, error)
	GetMagnetUri(ctx context.Context, in *GetMagnetUriReq, opts ...grpc.CallOption) (*GetMagnetUriRsp, error)
	GetResumeData(ctx context.Context, in *GetResumeDataReq, opts ...grpc.CallOption) (*GetResumeDataRsp, error)
	GetTorrentInfo(ctx context.Context, in *GetTorrentInfoReq, opts ...grpc.CallOption) (*GetTorrentInfoRsp, error)
	GetBtStatus(ctx context.Context, in *GetBtStatusReq, opts ...grpc.CallOption) (*GetBtStatusRsp, error)
	GetSessionParams(ctx context.Context, in *GetSessionParamsReq, opts ...grpc.CallOption) (*GetSessionParamsRsp, error)
	GetPeerInfo(ctx context.Context, in *GetPeerInfoReq, opts ...grpc.CallOption) (*GetPeerInfoRsp, error)
	OnBtStatus(ctx context.Context, opts ...grpc.CallOption) (BtService_OnBtStatusClient, error)
	OnFileCompleted(ctx context.Context, opts ...grpc.CallOption) (BtService_OnFileCompletedClient, error)
}

type btServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBtServiceClient(cc grpc.ClientConnInterface) BtServiceClient {
	return &btServiceClient{cc}
}

func (c *btServiceClient) InitedSession(ctx context.Context, in *InitedSessionReq, opts ...grpc.CallOption) (*InitedSessionRsp, error) {
	out := new(InitedSessionRsp)
	err := c.cc.Invoke(ctx, BtService_InitedSession_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) InitSession(ctx context.Context, in *InitSessionReq, opts ...grpc.CallOption) (*InitSessionRsp, error) {
	out := new(InitSessionRsp)
	err := c.cc.Invoke(ctx, BtService_InitSession_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) Parse(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadRespone, error) {
	out := new(DownloadRespone)
	err := c.cc.Invoke(ctx, BtService_Parse_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadRespone, error) {
	out := new(DownloadRespone)
	err := c.cc.Invoke(ctx, BtService_Download_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) RemoveTorrent(ctx context.Context, in *RemoveTorrentReq, opts ...grpc.CallOption) (*RemoveTorrentRes, error) {
	out := new(RemoveTorrentRes)
	err := c.cc.Invoke(ctx, BtService_RemoveTorrent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) GetMagnetUri(ctx context.Context, in *GetMagnetUriReq, opts ...grpc.CallOption) (*GetMagnetUriRsp, error) {
	out := new(GetMagnetUriRsp)
	err := c.cc.Invoke(ctx, BtService_GetMagnetUri_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) GetResumeData(ctx context.Context, in *GetResumeDataReq, opts ...grpc.CallOption) (*GetResumeDataRsp, error) {
	out := new(GetResumeDataRsp)
	err := c.cc.Invoke(ctx, BtService_GetResumeData_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) GetTorrentInfo(ctx context.Context, in *GetTorrentInfoReq, opts ...grpc.CallOption) (*GetTorrentInfoRsp, error) {
	out := new(GetTorrentInfoRsp)
	err := c.cc.Invoke(ctx, BtService_GetTorrentInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) GetBtStatus(ctx context.Context, in *GetBtStatusReq, opts ...grpc.CallOption) (*GetBtStatusRsp, error) {
	out := new(GetBtStatusRsp)
	err := c.cc.Invoke(ctx, BtService_GetBtStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) GetSessionParams(ctx context.Context, in *GetSessionParamsReq, opts ...grpc.CallOption) (*GetSessionParamsRsp, error) {
	out := new(GetSessionParamsRsp)
	err := c.cc.Invoke(ctx, BtService_GetSessionParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) GetPeerInfo(ctx context.Context, in *GetPeerInfoReq, opts ...grpc.CallOption) (*GetPeerInfoRsp, error) {
	out := new(GetPeerInfoRsp)
	err := c.cc.Invoke(ctx, BtService_GetPeerInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *btServiceClient) OnBtStatus(ctx context.Context, opts ...grpc.CallOption) (BtService_OnBtStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &BtService_ServiceDesc.Streams[0], BtService_OnBtStatus_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &btServiceOnBtStatusClient{stream}
	return x, nil
}

type BtService_OnBtStatusClient interface {
	Send(*BtStatusRequest) error
	Recv() (*BtStatusRespone, error)
	grpc.ClientStream
}

type btServiceOnBtStatusClient struct {
	grpc.ClientStream
}

func (x *btServiceOnBtStatusClient) Send(m *BtStatusRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *btServiceOnBtStatusClient) Recv() (*BtStatusRespone, error) {
	m := new(BtStatusRespone)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *btServiceClient) OnFileCompleted(ctx context.Context, opts ...grpc.CallOption) (BtService_OnFileCompletedClient, error) {
	stream, err := c.cc.NewStream(ctx, &BtService_ServiceDesc.Streams[1], BtService_OnFileCompleted_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &btServiceOnFileCompletedClient{stream}
	return x, nil
}

type BtService_OnFileCompletedClient interface {
	Send(*FileCompletedReq) error
	Recv() (*FileCompletedRes, error)
	grpc.ClientStream
}

type btServiceOnFileCompletedClient struct {
	grpc.ClientStream
}

func (x *btServiceOnFileCompletedClient) Send(m *FileCompletedReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *btServiceOnFileCompletedClient) Recv() (*FileCompletedRes, error) {
	m := new(FileCompletedRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BtServiceServer is the server API for BtService service.
// All implementations must embed UnimplementedBtServiceServer
// for forward compatibility
type BtServiceServer interface {
	InitedSession(context.Context, *InitedSessionReq) (*InitedSessionRsp, error)
	InitSession(context.Context, *InitSessionReq) (*InitSessionRsp, error)
	Parse(context.Context, *DownloadRequest) (*DownloadRespone, error)
	Download(context.Context, *DownloadRequest) (*DownloadRespone, error)
	RemoveTorrent(context.Context, *RemoveTorrentReq) (*RemoveTorrentRes, error)
	GetMagnetUri(context.Context, *GetMagnetUriReq) (*GetMagnetUriRsp, error)
	GetResumeData(context.Context, *GetResumeDataReq) (*GetResumeDataRsp, error)
	GetTorrentInfo(context.Context, *GetTorrentInfoReq) (*GetTorrentInfoRsp, error)
	GetBtStatus(context.Context, *GetBtStatusReq) (*GetBtStatusRsp, error)
	GetSessionParams(context.Context, *GetSessionParamsReq) (*GetSessionParamsRsp, error)
	GetPeerInfo(context.Context, *GetPeerInfoReq) (*GetPeerInfoRsp, error)
	OnBtStatus(BtService_OnBtStatusServer) error
	OnFileCompleted(BtService_OnFileCompletedServer) error
	mustEmbedUnimplementedBtServiceServer()
}

// UnimplementedBtServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBtServiceServer struct {
}

func (UnimplementedBtServiceServer) InitedSession(context.Context, *InitedSessionReq) (*InitedSessionRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitedSession not implemented")
}
func (UnimplementedBtServiceServer) InitSession(context.Context, *InitSessionReq) (*InitSessionRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitSession not implemented")
}
func (UnimplementedBtServiceServer) Parse(context.Context, *DownloadRequest) (*DownloadRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Parse not implemented")
}
func (UnimplementedBtServiceServer) Download(context.Context, *DownloadRequest) (*DownloadRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedBtServiceServer) RemoveTorrent(context.Context, *RemoveTorrentReq) (*RemoveTorrentRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveTorrent not implemented")
}
func (UnimplementedBtServiceServer) GetMagnetUri(context.Context, *GetMagnetUriReq) (*GetMagnetUriRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMagnetUri not implemented")
}
func (UnimplementedBtServiceServer) GetResumeData(context.Context, *GetResumeDataReq) (*GetResumeDataRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetResumeData not implemented")
}
func (UnimplementedBtServiceServer) GetTorrentInfo(context.Context, *GetTorrentInfoReq) (*GetTorrentInfoRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTorrentInfo not implemented")
}
func (UnimplementedBtServiceServer) GetBtStatus(context.Context, *GetBtStatusReq) (*GetBtStatusRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBtStatus not implemented")
}
func (UnimplementedBtServiceServer) GetSessionParams(context.Context, *GetSessionParamsReq) (*GetSessionParamsRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSessionParams not implemented")
}
func (UnimplementedBtServiceServer) GetPeerInfo(context.Context, *GetPeerInfoReq) (*GetPeerInfoRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPeerInfo not implemented")
}
func (UnimplementedBtServiceServer) OnBtStatus(BtService_OnBtStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method OnBtStatus not implemented")
}
func (UnimplementedBtServiceServer) OnFileCompleted(BtService_OnFileCompletedServer) error {
	return status.Errorf(codes.Unimplemented, "method OnFileCompleted not implemented")
}
func (UnimplementedBtServiceServer) mustEmbedUnimplementedBtServiceServer() {}

// UnsafeBtServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BtServiceServer will
// result in compilation errors.
type UnsafeBtServiceServer interface {
	mustEmbedUnimplementedBtServiceServer()
}

func RegisterBtServiceServer(s grpc.ServiceRegistrar, srv BtServiceServer) {
	s.RegisterService(&BtService_ServiceDesc, srv)
}

func _BtService_InitedSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitedSessionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).InitedSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_InitedSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).InitedSession(ctx, req.(*InitedSessionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_InitSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitSessionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).InitSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_InitSession_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).InitSession(ctx, req.(*InitSessionReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_Parse_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).Parse(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_Parse_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).Parse(ctx, req.(*DownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_Download_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).Download(ctx, req.(*DownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_RemoveTorrent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveTorrentReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).RemoveTorrent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_RemoveTorrent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).RemoveTorrent(ctx, req.(*RemoveTorrentReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_GetMagnetUri_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMagnetUriReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).GetMagnetUri(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_GetMagnetUri_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).GetMagnetUri(ctx, req.(*GetMagnetUriReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_GetResumeData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetResumeDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).GetResumeData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_GetResumeData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).GetResumeData(ctx, req.(*GetResumeDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_GetTorrentInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTorrentInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).GetTorrentInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_GetTorrentInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).GetTorrentInfo(ctx, req.(*GetTorrentInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_GetBtStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBtStatusReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).GetBtStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_GetBtStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).GetBtStatus(ctx, req.(*GetBtStatusReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_GetSessionParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSessionParamsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).GetSessionParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_GetSessionParams_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).GetSessionParams(ctx, req.(*GetSessionParamsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_GetPeerInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPeerInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BtServiceServer).GetPeerInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BtService_GetPeerInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BtServiceServer).GetPeerInfo(ctx, req.(*GetPeerInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _BtService_OnBtStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BtServiceServer).OnBtStatus(&btServiceOnBtStatusServer{stream})
}

type BtService_OnBtStatusServer interface {
	Send(*BtStatusRespone) error
	Recv() (*BtStatusRequest, error)
	grpc.ServerStream
}

type btServiceOnBtStatusServer struct {
	grpc.ServerStream
}

func (x *btServiceOnBtStatusServer) Send(m *BtStatusRespone) error {
	return x.ServerStream.SendMsg(m)
}

func (x *btServiceOnBtStatusServer) Recv() (*BtStatusRequest, error) {
	m := new(BtStatusRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _BtService_OnFileCompleted_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BtServiceServer).OnFileCompleted(&btServiceOnFileCompletedServer{stream})
}

type BtService_OnFileCompletedServer interface {
	Send(*FileCompletedRes) error
	Recv() (*FileCompletedReq, error)
	grpc.ServerStream
}

type btServiceOnFileCompletedServer struct {
	grpc.ServerStream
}

func (x *btServiceOnFileCompletedServer) Send(m *FileCompletedRes) error {
	return x.ServerStream.SendMsg(m)
}

func (x *btServiceOnFileCompletedServer) Recv() (*FileCompletedReq, error) {
	m := new(FileCompletedReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BtService_ServiceDesc is the grpc.ServiceDesc for BtService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BtService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "prpc.BtService",
	HandlerType: (*BtServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitedSession",
			Handler:    _BtService_InitedSession_Handler,
		},
		{
			MethodName: "InitSession",
			Handler:    _BtService_InitSession_Handler,
		},
		{
			MethodName: "Parse",
			Handler:    _BtService_Parse_Handler,
		},
		{
			MethodName: "Download",
			Handler:    _BtService_Download_Handler,
		},
		{
			MethodName: "RemoveTorrent",
			Handler:    _BtService_RemoveTorrent_Handler,
		},
		{
			MethodName: "GetMagnetUri",
			Handler:    _BtService_GetMagnetUri_Handler,
		},
		{
			MethodName: "GetResumeData",
			Handler:    _BtService_GetResumeData_Handler,
		},
		{
			MethodName: "GetTorrentInfo",
			Handler:    _BtService_GetTorrentInfo_Handler,
		},
		{
			MethodName: "GetBtStatus",
			Handler:    _BtService_GetBtStatus_Handler,
		},
		{
			MethodName: "GetSessionParams",
			Handler:    _BtService_GetSessionParams_Handler,
		},
		{
			MethodName: "GetPeerInfo",
			Handler:    _BtService_GetPeerInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "OnBtStatus",
			Handler:       _BtService_OnBtStatus_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "OnFileCompleted",
			Handler:       _BtService_OnFileCompleted_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "bt.proto",
}
