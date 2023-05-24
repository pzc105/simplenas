/**
 * @fileoverview gRPC-Web generated client stub for prpc
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.4.2
// 	protoc              v3.19.1
// source: user.proto


/* eslint-disable */
// @ts-nocheck


import * as grpcWeb from 'grpc-web';

import * as user_pb from './user_pb';
import * as bt_pb from './bt_pb';


export class UserServiceClient {
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

  methodDescriptorRegister = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/Register',
    grpcWeb.MethodType.UNARY,
    user_pb.RegisterInfo,
    user_pb.RegisterRet,
    (request: user_pb.RegisterInfo) => {
      return request.serializeBinary();
    },
    user_pb.RegisterRet.deserializeBinary
  );

  register(
    request: user_pb.RegisterInfo,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.RegisterRet>;

  register(
    request: user_pb.RegisterInfo,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.RegisterRet) => void): grpcWeb.ClientReadableStream<user_pb.RegisterRet>;

  register(
    request: user_pb.RegisterInfo,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.RegisterRet) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/Register',
        request,
        metadata || {},
        this.methodDescriptorRegister,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/Register',
    request,
    metadata || {},
    this.methodDescriptorRegister);
  }

  methodDescriptorIsUsedEmail = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/IsUsedEmail',
    grpcWeb.MethodType.UNARY,
    user_pb.EmailInfo,
    user_pb.IsUsedEmailRet,
    (request: user_pb.EmailInfo) => {
      return request.serializeBinary();
    },
    user_pb.IsUsedEmailRet.deserializeBinary
  );

  isUsedEmail(
    request: user_pb.EmailInfo,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.IsUsedEmailRet>;

  isUsedEmail(
    request: user_pb.EmailInfo,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.IsUsedEmailRet) => void): grpcWeb.ClientReadableStream<user_pb.IsUsedEmailRet>;

  isUsedEmail(
    request: user_pb.EmailInfo,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.IsUsedEmailRet) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/IsUsedEmail',
        request,
        metadata || {},
        this.methodDescriptorIsUsedEmail,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/IsUsedEmail',
    request,
    metadata || {},
    this.methodDescriptorIsUsedEmail);
  }

  methodDescriptorLogin = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/Login',
    grpcWeb.MethodType.UNARY,
    user_pb.LoginInfo,
    user_pb.LoginRet,
    (request: user_pb.LoginInfo) => {
      return request.serializeBinary();
    },
    user_pb.LoginRet.deserializeBinary
  );

  login(
    request: user_pb.LoginInfo,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.LoginRet>;

  login(
    request: user_pb.LoginInfo,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.LoginRet) => void): grpcWeb.ClientReadableStream<user_pb.LoginRet>;

  login(
    request: user_pb.LoginInfo,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.LoginRet) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/Login',
        request,
        metadata || {},
        this.methodDescriptorLogin,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/Login',
    request,
    metadata || {},
    this.methodDescriptorLogin);
  }

  methodDescriptorFastLogin = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/FastLogin',
    grpcWeb.MethodType.UNARY,
    user_pb.LoginInfo,
    user_pb.LoginRet,
    (request: user_pb.LoginInfo) => {
      return request.serializeBinary();
    },
    user_pb.LoginRet.deserializeBinary
  );

  fastLogin(
    request: user_pb.LoginInfo,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.LoginRet>;

  fastLogin(
    request: user_pb.LoginInfo,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.LoginRet) => void): grpcWeb.ClientReadableStream<user_pb.LoginRet>;

  fastLogin(
    request: user_pb.LoginInfo,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.LoginRet) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/FastLogin',
        request,
        metadata || {},
        this.methodDescriptorFastLogin,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/FastLogin',
    request,
    metadata || {},
    this.methodDescriptorFastLogin);
  }

  methodDescriptorIsLogined = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/IsLogined',
    grpcWeb.MethodType.UNARY,
    user_pb.LoginInfo,
    user_pb.LoginRet,
    (request: user_pb.LoginInfo) => {
      return request.serializeBinary();
    },
    user_pb.LoginRet.deserializeBinary
  );

  isLogined(
    request: user_pb.LoginInfo,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.LoginRet>;

  isLogined(
    request: user_pb.LoginInfo,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.LoginRet) => void): grpcWeb.ClientReadableStream<user_pb.LoginRet>;

  isLogined(
    request: user_pb.LoginInfo,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.LoginRet) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/IsLogined',
        request,
        metadata || {},
        this.methodDescriptorIsLogined,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/IsLogined',
    request,
    metadata || {},
    this.methodDescriptorIsLogined);
  }

  methodDescriptorDownload = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/Download',
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
          '/prpc.UserService/Download',
        request,
        metadata || {},
        this.methodDescriptorDownload,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/Download',
    request,
    metadata || {},
    this.methodDescriptorDownload);
  }

  methodDescriptorRemoveTorrent = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/RemoveTorrent',
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
          '/prpc.UserService/RemoveTorrent',
        request,
        metadata || {},
        this.methodDescriptorRemoveTorrent,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/RemoveTorrent',
    request,
    metadata || {},
    this.methodDescriptorRemoveTorrent);
  }

  methodDescriptorOnStatus = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/OnStatus',
    grpcWeb.MethodType.SERVER_STREAMING,
    bt_pb.StatusRequest,
    bt_pb.StatusRespone,
    (request: bt_pb.StatusRequest) => {
      return request.serializeBinary();
    },
    bt_pb.StatusRespone.deserializeBinary
  );

  onStatus(
    request: bt_pb.StatusRequest,
    metadata?: grpcWeb.Metadata): grpcWeb.ClientReadableStream<bt_pb.StatusRespone> {
    return this.client_.serverStreaming(
      this.hostname_ +
        '/prpc.UserService/OnStatus',
      request,
      metadata || {},
      this.methodDescriptorOnStatus);
  }

  methodDescriptorQueryBtVideos = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/QueryBtVideos',
    grpcWeb.MethodType.UNARY,
    user_pb.QueryBtVideosReq,
    user_pb.QueryBtVideosRes,
    (request: user_pb.QueryBtVideosReq) => {
      return request.serializeBinary();
    },
    user_pb.QueryBtVideosRes.deserializeBinary
  );

  queryBtVideos(
    request: user_pb.QueryBtVideosReq,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.QueryBtVideosRes>;

  queryBtVideos(
    request: user_pb.QueryBtVideosReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.QueryBtVideosRes) => void): grpcWeb.ClientReadableStream<user_pb.QueryBtVideosRes>;

  queryBtVideos(
    request: user_pb.QueryBtVideosReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.QueryBtVideosRes) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/QueryBtVideos',
        request,
        metadata || {},
        this.methodDescriptorQueryBtVideos,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/QueryBtVideos',
    request,
    metadata || {},
    this.methodDescriptorQueryBtVideos);
  }

  methodDescriptorNewCategoryItem = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/NewCategoryItem',
    grpcWeb.MethodType.UNARY,
    user_pb.NewCategoryItemReq,
    user_pb.NewCategoryItemRes,
    (request: user_pb.NewCategoryItemReq) => {
      return request.serializeBinary();
    },
    user_pb.NewCategoryItemRes.deserializeBinary
  );

  newCategoryItem(
    request: user_pb.NewCategoryItemReq,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.NewCategoryItemRes>;

  newCategoryItem(
    request: user_pb.NewCategoryItemReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.NewCategoryItemRes) => void): grpcWeb.ClientReadableStream<user_pb.NewCategoryItemRes>;

  newCategoryItem(
    request: user_pb.NewCategoryItemReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.NewCategoryItemRes) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/NewCategoryItem',
        request,
        metadata || {},
        this.methodDescriptorNewCategoryItem,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/NewCategoryItem',
    request,
    metadata || {},
    this.methodDescriptorNewCategoryItem);
  }

  methodDescriptorDelCategoryItem = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/DelCategoryItem',
    grpcWeb.MethodType.UNARY,
    user_pb.DelCategoryItemReq,
    user_pb.DelCategoryItemRes,
    (request: user_pb.DelCategoryItemReq) => {
      return request.serializeBinary();
    },
    user_pb.DelCategoryItemRes.deserializeBinary
  );

  delCategoryItem(
    request: user_pb.DelCategoryItemReq,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.DelCategoryItemRes>;

  delCategoryItem(
    request: user_pb.DelCategoryItemReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.DelCategoryItemRes) => void): grpcWeb.ClientReadableStream<user_pb.DelCategoryItemRes>;

  delCategoryItem(
    request: user_pb.DelCategoryItemReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.DelCategoryItemRes) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/DelCategoryItem',
        request,
        metadata || {},
        this.methodDescriptorDelCategoryItem,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/DelCategoryItem',
    request,
    metadata || {},
    this.methodDescriptorDelCategoryItem);
  }

  methodDescriptorAddBtVideos = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/AddBtVideos',
    grpcWeb.MethodType.UNARY,
    user_pb.AddBtVideosReq,
    user_pb.AddBtVideosRes,
    (request: user_pb.AddBtVideosReq) => {
      return request.serializeBinary();
    },
    user_pb.AddBtVideosRes.deserializeBinary
  );

  addBtVideos(
    request: user_pb.AddBtVideosReq,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.AddBtVideosRes>;

  addBtVideos(
    request: user_pb.AddBtVideosReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.AddBtVideosRes) => void): grpcWeb.ClientReadableStream<user_pb.AddBtVideosRes>;

  addBtVideos(
    request: user_pb.AddBtVideosReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.AddBtVideosRes) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/AddBtVideos',
        request,
        metadata || {},
        this.methodDescriptorAddBtVideos,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/AddBtVideos',
    request,
    metadata || {},
    this.methodDescriptorAddBtVideos);
  }

  methodDescriptorShareItem = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/ShareItem',
    grpcWeb.MethodType.UNARY,
    user_pb.ShareItemReq,
    user_pb.ShareItemRes,
    (request: user_pb.ShareItemReq) => {
      return request.serializeBinary();
    },
    user_pb.ShareItemRes.deserializeBinary
  );

  shareItem(
    request: user_pb.ShareItemReq,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.ShareItemRes>;

  shareItem(
    request: user_pb.ShareItemReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.ShareItemRes) => void): grpcWeb.ClientReadableStream<user_pb.ShareItemRes>;

  shareItem(
    request: user_pb.ShareItemReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.ShareItemRes) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/ShareItem',
        request,
        metadata || {},
        this.methodDescriptorShareItem,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/ShareItem',
    request,
    metadata || {},
    this.methodDescriptorShareItem);
  }

  methodDescriptorQuerySubItems = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/QuerySubItems',
    grpcWeb.MethodType.UNARY,
    user_pb.QuerySubItemsReq,
    user_pb.QuerySubItemsRes,
    (request: user_pb.QuerySubItemsReq) => {
      return request.serializeBinary();
    },
    user_pb.QuerySubItemsRes.deserializeBinary
  );

  querySubItems(
    request: user_pb.QuerySubItemsReq,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.QuerySubItemsRes>;

  querySubItems(
    request: user_pb.QuerySubItemsReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.QuerySubItemsRes) => void): grpcWeb.ClientReadableStream<user_pb.QuerySubItemsRes>;

  querySubItems(
    request: user_pb.QuerySubItemsReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.QuerySubItemsRes) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/QuerySubItems',
        request,
        metadata || {},
        this.methodDescriptorQuerySubItems,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/QuerySubItems',
    request,
    metadata || {},
    this.methodDescriptorQuerySubItems);
  }

  methodDescriptorQueryItemInfo = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/QueryItemInfo',
    grpcWeb.MethodType.UNARY,
    user_pb.QueryItemInfoReq,
    user_pb.QueryItemInfoRes,
    (request: user_pb.QueryItemInfoReq) => {
      return request.serializeBinary();
    },
    user_pb.QueryItemInfoRes.deserializeBinary
  );

  queryItemInfo(
    request: user_pb.QueryItemInfoReq,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.QueryItemInfoRes>;

  queryItemInfo(
    request: user_pb.QueryItemInfoReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.QueryItemInfoRes) => void): grpcWeb.ClientReadableStream<user_pb.QueryItemInfoRes>;

  queryItemInfo(
    request: user_pb.QueryItemInfoReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.QueryItemInfoRes) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/QueryItemInfo',
        request,
        metadata || {},
        this.methodDescriptorQueryItemInfo,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/QueryItemInfo',
    request,
    metadata || {},
    this.methodDescriptorQueryItemInfo);
  }

  methodDescriptorRefreshSubtitle = new grpcWeb.MethodDescriptor(
    '/prpc.UserService/RefreshSubtitle',
    grpcWeb.MethodType.UNARY,
    user_pb.RefreshSubtitleReq,
    user_pb.RefreshSubtitleRes,
    (request: user_pb.RefreshSubtitleReq) => {
      return request.serializeBinary();
    },
    user_pb.RefreshSubtitleRes.deserializeBinary
  );

  refreshSubtitle(
    request: user_pb.RefreshSubtitleReq,
    metadata: grpcWeb.Metadata | null): Promise<user_pb.RefreshSubtitleRes>;

  refreshSubtitle(
    request: user_pb.RefreshSubtitleReq,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.RpcError,
               response: user_pb.RefreshSubtitleRes) => void): grpcWeb.ClientReadableStream<user_pb.RefreshSubtitleRes>;

  refreshSubtitle(
    request: user_pb.RefreshSubtitleReq,
    metadata: grpcWeb.Metadata | null,
    callback?: (err: grpcWeb.RpcError,
               response: user_pb.RefreshSubtitleRes) => void) {
    if (callback !== undefined) {
      return this.client_.rpcCall(
        this.hostname_ +
          '/prpc.UserService/RefreshSubtitle',
        request,
        metadata || {},
        this.methodDescriptorRefreshSubtitle,
        callback);
    }
    return this.client_.unaryCall(
    this.hostname_ +
      '/prpc.UserService/RefreshSubtitle',
    request,
    metadata || {},
    this.methodDescriptorRefreshSubtitle);
  }

}

