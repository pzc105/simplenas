// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
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
	BtService_Parse_FullMethodName           = "/prpc.BtService/Parse"
	BtService_Download_FullMethodName        = "/prpc.BtService/Download"
	BtService_RemoveTorrent_FullMethodName   = "/prpc.BtService/RemoveTorrent"
	BtService_OnStatus_FullMethodName        = "/prpc.BtService/OnStatus"
	BtService_OnTorrentInfo_FullMethodName   = "/prpc.BtService/OnTorrentInfo"
	BtService_OnFileCompleted_FullMethodName = "/prpc.BtService/OnFileCompleted"
	BtService_FileProgress_FullMethodName    = "/prpc.BtService/FileProgress"
)

// BtServiceClient is the client API for BtService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BtServiceClient interface {
	Parse(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadRespone, error)
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadRespone, error)
	RemoveTorrent(ctx context.Context, in *RemoveTorrentReq, opts ...grpc.CallOption) (*RemoveTorrentRes, error)
	OnStatus(ctx context.Context, opts ...grpc.CallOption) (BtService_OnStatusClient, error)
	OnTorrentInfo(ctx context.Context, opts ...grpc.CallOption) (BtService_OnTorrentInfoClient, error)
	OnFileCompleted(ctx context.Context, opts ...grpc.CallOption) (BtService_OnFileCompletedClient, error)
	FileProgress(ctx context.Context, opts ...grpc.CallOption) (BtService_FileProgressClient, error)
}

type btServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBtServiceClient(cc grpc.ClientConnInterface) BtServiceClient {
	return &btServiceClient{cc}
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

func (c *btServiceClient) OnStatus(ctx context.Context, opts ...grpc.CallOption) (BtService_OnStatusClient, error) {
	stream, err := c.cc.NewStream(ctx, &BtService_ServiceDesc.Streams[0], BtService_OnStatus_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &btServiceOnStatusClient{stream}
	return x, nil
}

type BtService_OnStatusClient interface {
	Send(*StatusRequest) error
	Recv() (*StatusRespone, error)
	grpc.ClientStream
}

type btServiceOnStatusClient struct {
	grpc.ClientStream
}

func (x *btServiceOnStatusClient) Send(m *StatusRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *btServiceOnStatusClient) Recv() (*StatusRespone, error) {
	m := new(StatusRespone)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *btServiceClient) OnTorrentInfo(ctx context.Context, opts ...grpc.CallOption) (BtService_OnTorrentInfoClient, error) {
	stream, err := c.cc.NewStream(ctx, &BtService_ServiceDesc.Streams[1], BtService_OnTorrentInfo_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &btServiceOnTorrentInfoClient{stream}
	return x, nil
}

type BtService_OnTorrentInfoClient interface {
	Send(*TorrentInfoReq) error
	Recv() (*TorrentInfoRes, error)
	grpc.ClientStream
}

type btServiceOnTorrentInfoClient struct {
	grpc.ClientStream
}

func (x *btServiceOnTorrentInfoClient) Send(m *TorrentInfoReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *btServiceOnTorrentInfoClient) Recv() (*TorrentInfoRes, error) {
	m := new(TorrentInfoRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *btServiceClient) OnFileCompleted(ctx context.Context, opts ...grpc.CallOption) (BtService_OnFileCompletedClient, error) {
	stream, err := c.cc.NewStream(ctx, &BtService_ServiceDesc.Streams[2], BtService_OnFileCompleted_FullMethodName, opts...)
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

func (c *btServiceClient) FileProgress(ctx context.Context, opts ...grpc.CallOption) (BtService_FileProgressClient, error) {
	stream, err := c.cc.NewStream(ctx, &BtService_ServiceDesc.Streams[3], BtService_FileProgress_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &btServiceFileProgressClient{stream}
	return x, nil
}

type BtService_FileProgressClient interface {
	Send(*FileProgressReq) error
	Recv() (*FileProgressRes, error)
	grpc.ClientStream
}

type btServiceFileProgressClient struct {
	grpc.ClientStream
}

func (x *btServiceFileProgressClient) Send(m *FileProgressReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *btServiceFileProgressClient) Recv() (*FileProgressRes, error) {
	m := new(FileProgressRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BtServiceServer is the server API for BtService service.
// All implementations must embed UnimplementedBtServiceServer
// for forward compatibility
type BtServiceServer interface {
	Parse(context.Context, *DownloadRequest) (*DownloadRespone, error)
	Download(context.Context, *DownloadRequest) (*DownloadRespone, error)
	RemoveTorrent(context.Context, *RemoveTorrentReq) (*RemoveTorrentRes, error)
	OnStatus(BtService_OnStatusServer) error
	OnTorrentInfo(BtService_OnTorrentInfoServer) error
	OnFileCompleted(BtService_OnFileCompletedServer) error
	FileProgress(BtService_FileProgressServer) error
	mustEmbedUnimplementedBtServiceServer()
}

// UnimplementedBtServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBtServiceServer struct {
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
func (UnimplementedBtServiceServer) OnStatus(BtService_OnStatusServer) error {
	return status.Errorf(codes.Unimplemented, "method OnStatus not implemented")
}
func (UnimplementedBtServiceServer) OnTorrentInfo(BtService_OnTorrentInfoServer) error {
	return status.Errorf(codes.Unimplemented, "method OnTorrentInfo not implemented")
}
func (UnimplementedBtServiceServer) OnFileCompleted(BtService_OnFileCompletedServer) error {
	return status.Errorf(codes.Unimplemented, "method OnFileCompleted not implemented")
}
func (UnimplementedBtServiceServer) FileProgress(BtService_FileProgressServer) error {
	return status.Errorf(codes.Unimplemented, "method FileProgress not implemented")
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

func _BtService_OnStatus_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BtServiceServer).OnStatus(&btServiceOnStatusServer{stream})
}

type BtService_OnStatusServer interface {
	Send(*StatusRespone) error
	Recv() (*StatusRequest, error)
	grpc.ServerStream
}

type btServiceOnStatusServer struct {
	grpc.ServerStream
}

func (x *btServiceOnStatusServer) Send(m *StatusRespone) error {
	return x.ServerStream.SendMsg(m)
}

func (x *btServiceOnStatusServer) Recv() (*StatusRequest, error) {
	m := new(StatusRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _BtService_OnTorrentInfo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BtServiceServer).OnTorrentInfo(&btServiceOnTorrentInfoServer{stream})
}

type BtService_OnTorrentInfoServer interface {
	Send(*TorrentInfoRes) error
	Recv() (*TorrentInfoReq, error)
	grpc.ServerStream
}

type btServiceOnTorrentInfoServer struct {
	grpc.ServerStream
}

func (x *btServiceOnTorrentInfoServer) Send(m *TorrentInfoRes) error {
	return x.ServerStream.SendMsg(m)
}

func (x *btServiceOnTorrentInfoServer) Recv() (*TorrentInfoReq, error) {
	m := new(TorrentInfoReq)
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

func _BtService_FileProgress_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BtServiceServer).FileProgress(&btServiceFileProgressServer{stream})
}

type BtService_FileProgressServer interface {
	Send(*FileProgressRes) error
	Recv() (*FileProgressReq, error)
	grpc.ServerStream
}

type btServiceFileProgressServer struct {
	grpc.ServerStream
}

func (x *btServiceFileProgressServer) Send(m *FileProgressRes) error {
	return x.ServerStream.SendMsg(m)
}

func (x *btServiceFileProgressServer) Recv() (*FileProgressReq, error) {
	m := new(FileProgressReq)
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
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "OnStatus",
			Handler:       _BtService_OnStatus_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "OnTorrentInfo",
			Handler:       _BtService_OnTorrentInfo_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "OnFileCompleted",
			Handler:       _BtService_OnFileCompleted_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "FileProgress",
			Handler:       _BtService_FileProgress_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "bt.proto",
}
