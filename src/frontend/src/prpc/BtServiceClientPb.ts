/**
 * @fileoverview gRPC-Web generated client stub for prpc
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.4.2
// 	protoc              v4.24.3
// source: bt.proto


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as bt_pb from './bt_pb';


export class BtServiceClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: any; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname.replace(/\/+$/, '');
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodDescriptorParse = new grpcWeb.MethodDescriptor(
    '/prpc.BtService/Parse',
    grpcWeb.MethodType.UNARY,
    bt_pb.DownloadRequest,
    bt_pb.DownloadRespone,
    (request: bt_pb.DownloadRequest) => {
      return request.serializeBinary();
    },
    bt_pb.DownloadRespone.deserializeBinary
  );

  parse(
    request: bt_pb.DownloadRequest,
    metadata: grpcWeb.Metadata | null): Promise<bt_pb.DownloadRespone>;

  parse(
    request: bt_pb.DownloadRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: bt_pb.DownloadRespone) => void): grpcWeb.ClientReadableStream<bt_pb.DownloadRespone>;

  parse(
    request: bt_pb.DownloadRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: bt_pb.DownloadRespone) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.BtService/Parse',
        request,
        metadata || {},
        this.methodDescriptorParse,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.BtService/Parse',
    request,
    metadata || {},
    this.methodDescriptorParse);
  }

  methodDescriptorDownload = new grpcWeb.MethodDescriptor(
    '/prpc.BtService/Download',
    grpcWeb.MethodType.UNARY,
    bt_pb.DownloadRequest,
    bt_pb.DownloadRespone,
    (request: bt_pb.DownloadRequest) => {
      return request.serializeBinary();
    },
    bt_pb.DownloadRespone.deserializeBinary
  );

  download(
    request: bt_pb.DownloadRequest,
    metadata: grpcWeb.Metadata | null): Promise<bt_pb.DownloadRespone>;

  download(
    request: bt_pb.DownloadRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: bt_pb.DownloadRespone) => void): grpcWeb.ClientReadableStream<bt_pb.DownloadRespone>;

  download(
    request: bt_pb.DownloadRequest,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: bt_pb.DownloadRespone) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.BtService/Download',
        request,
        metadata || {},
        this.methodDescriptorDownload,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.BtService/Download',
    request,
    metadata || {},
    this.methodDescriptorDownload);
  }

  methodDescriptorRemoveTorrent = new grpcWeb.MethodDescriptor(
    '/prpc.BtService/RemoveTorrent',
    grpcWeb.MethodType.UNARY,
    bt_pb.RemoveTorrentReq,
    bt_pb.RemoveTorrentRes,
    (request: bt_pb.RemoveTorrentReq) => {
      return request.serializeBinary();
    },
    bt_pb.RemoveTorrentRes.deserializeBinary
  );

  removeTorrent(
    request: bt_pb.RemoveTorrentReq,
    metadata: grpcWeb.Metadata | null): Promise<bt_pb.RemoveTorrentRes>;

  removeTorrent(
    request: bt_pb.RemoveTorrentReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: bt_pb.RemoveTorrentRes) => void): grpcWeb.ClientReadableStream<bt_pb.RemoveTorrentRes>;

  removeTorrent(
    request: bt_pb.RemoveTorrentReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: bt_pb.RemoveTorrentRes) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.BtService/RemoveTorrent',
        request,
        metadata || {},
        this.methodDescriptorRemoveTorrent,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.BtService/RemoveTorrent',
    request,
    metadata || {},
    this.methodDescriptorRemoveTorrent);
  }

  methodDescriptorGetMagnetUri = new grpcWeb.MethodDescriptor(
    '/prpc.BtService/GetMagnetUri',
    grpcWeb.MethodType.UNARY,
    bt_pb.GetMagnetUriReq,
    bt_pb.GetMagnetUriRsp,
    (request: bt_pb.GetMagnetUriReq) => {
      return request.serializeBinary();
    },
    bt_pb.GetMagnetUriRsp.deserializeBinary
  );

  getMagnetUri(
    request: bt_pb.GetMagnetUriReq,
    metadata: grpcWeb.Metadata | null): Promise<bt_pb.GetMagnetUriRsp>;

  getMagnetUri(
    request: bt_pb.GetMagnetUriReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: bt_pb.GetMagnetUriRsp) => void): grpcWeb.ClientReadableStream<bt_pb.GetMagnetUriRsp>;

  getMagnetUri(
    request: bt_pb.GetMagnetUriReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: bt_pb.GetMagnetUriRsp) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.BtService/GetMagnetUri',
        request,
        metadata || {},
        this.methodDescriptorGetMagnetUri,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.BtService/GetMagnetUri',
    request,
    metadata || {},
    this.methodDescriptorGetMagnetUri);
  }

  methodDescriptorGetResumeData = new grpcWeb.MethodDescriptor(
    '/prpc.BtService/GetResumeData',
    grpcWeb.MethodType.UNARY,
    bt_pb.GetResumeDataReq,
    bt_pb.GetResumeDataRsp,
    (request: bt_pb.GetResumeDataReq) => {
      return request.serializeBinary();
    },
    bt_pb.GetResumeDataRsp.deserializeBinary
  );

  getResumeData(
    request: bt_pb.GetResumeDataReq,
    metadata: grpcWeb.Metadata | null): Promise<bt_pb.GetResumeDataRsp>;

  getResumeData(
    request: bt_pb.GetResumeDataReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: bt_pb.GetResumeDataRsp) => void): grpcWeb.ClientReadableStream<bt_pb.GetResumeDataRsp>;

  getResumeData(
    request: bt_pb.GetResumeDataReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: bt_pb.GetResumeDataRsp) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.BtService/GetResumeData',
        request,
        metadata || {},
        this.methodDescriptorGetResumeData,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.BtService/GetResumeData',
    request,
    metadata || {},
    this.methodDescriptorGetResumeData);
  }

  methodDescriptorGetTorrentInfo = new grpcWeb.MethodDescriptor(
    '/prpc.BtService/GetTorrentInfo',
    grpcWeb.MethodType.UNARY,
    bt_pb.GetTorrentInfoReq,
    bt_pb.GetTorrentInfoRsp,
    (request: bt_pb.GetTorrentInfoReq) => {
      return request.serializeBinary();
    },
    bt_pb.GetTorrentInfoRsp.deserializeBinary
  );

  getTorrentInfo(
    request: bt_pb.GetTorrentInfoReq,
    metadata: grpcWeb.Metadata | null): Promise<bt_pb.GetTorrentInfoRsp>;

  getTorrentInfo(
    request: bt_pb.GetTorrentInfoReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: bt_pb.GetTorrentInfoRsp) => void): grpcWeb.ClientReadableStream<bt_pb.GetTorrentInfoRsp>;

  getTorrentInfo(
    request: bt_pb.GetTorrentInfoReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: bt_pb.GetTorrentInfoRsp) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.BtService/GetTorrentInfo',
        request,
        metadata || {},
        this.methodDescriptorGetTorrentInfo,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.BtService/GetTorrentInfo',
    request,
    metadata || {},
    this.methodDescriptorGetTorrentInfo);
  }

  methodDescriptorGetBtStatus = new grpcWeb.MethodDescriptor(
    '/prpc.BtService/GetBtStatus',
    grpcWeb.MethodType.UNARY,
    bt_pb.GetBtStatusReq,
    bt_pb.GetBtStatusRsp,
    (request: bt_pb.GetBtStatusReq) => {
      return request.serializeBinary();
    },
    bt_pb.GetBtStatusRsp.deserializeBinary
  );

  getBtStatus(
    request: bt_pb.GetBtStatusReq,
    metadata: grpcWeb.Metadata | null): Promise<bt_pb.GetBtStatusRsp>;

  getBtStatus(
    request: bt_pb.GetBtStatusReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: bt_pb.GetBtStatusRsp) => void): grpcWeb.ClientReadableStream<bt_pb.GetBtStatusRsp>;

  getBtStatus(
    request: bt_pb.GetBtStatusReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: bt_pb.GetBtStatusRsp) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.BtService/GetBtStatus',
        request,
        metadata || {},
        this.methodDescriptorGetBtStatus,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.BtService/GetBtStatus',
    request,
    metadata || {},
    this.methodDescriptorGetBtStatus);
  }

}

