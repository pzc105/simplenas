import * as jspb from 'google-protobuf'



export class ChatError extends jspb.Message {
  getErrorId(): ChatError.ErrorId;
  setErrorId(value: ChatError.ErrorId): ChatError;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChatError.AsObject;
  static toObject(includeInstance: boolean, msg: ChatError): ChatError.AsObject;
  static serializeBinaryToWriter(message: ChatError, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChatError;
  static deserializeBinaryFromReader(message: ChatError, reader: jspb.BinaryReader): ChatError;
}

export namespace ChatError {
  export type AsObject = {
    errorId: ChatError.ErrorId,
  }

  export enum ErrorId { 
    NONE = 0,
    INVALIDIDENTITY = 1,
    NOTSUPPORTED = 2,
  }
}

export class ChatUserInfo extends jspb.Message {
  getUserId(): number;
  setUserId(value: number): ChatUserInfo;

  getToken(): string;
  setToken(value: string): ChatUserInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChatUserInfo.AsObject;
  static toObject(includeInstance: boolean, msg: ChatUserInfo): ChatUserInfo.AsObject;
  static serializeBinaryToWriter(message: ChatUserInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChatUserInfo;
  static deserializeBinaryFromReader(message: ChatUserInfo, reader: jspb.BinaryReader): ChatUserInfo;
}

export namespace ChatUserInfo {
  export type AsObject = {
    userId: number,
    token: string,
  }
}

export class QueryChatRoomServerReq extends jspb.Message {
  getRoomKey(): string;
  setRoomKey(value: string): QueryChatRoomServerReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryChatRoomServerReq.AsObject;
  static toObject(includeInstance: boolean, msg: QueryChatRoomServerReq): QueryChatRoomServerReq.AsObject;
  static serializeBinaryToWriter(message: QueryChatRoomServerReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryChatRoomServerReq;
  static deserializeBinaryFromReader(message: QueryChatRoomServerReq, reader: jspb.BinaryReader): QueryChatRoomServerReq;
}

export namespace QueryChatRoomServerReq {
  export type AsObject = {
    roomKey: string,
  }
}

export class QueryChatRoomServerRes extends jspb.Message {
  getAddressesList(): Array<string>;
  setAddressesList(value: Array<string>): QueryChatRoomServerRes;
  clearAddressesList(): QueryChatRoomServerRes;
  addAddresses(value: string, index?: number): QueryChatRoomServerRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryChatRoomServerRes.AsObject;
  static toObject(includeInstance: boolean, msg: QueryChatRoomServerRes): QueryChatRoomServerRes.AsObject;
  static serializeBinaryToWriter(message: QueryChatRoomServerRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryChatRoomServerRes;
  static deserializeBinaryFromReader(message: QueryChatRoomServerRes, reader: jspb.BinaryReader): QueryChatRoomServerRes;
}

export namespace QueryChatRoomServerRes {
  export type AsObject = {
    addressesList: Array<string>,
  }
}

export class CreateChatRoomReq extends jspb.Message {
  getUserInfo(): ChatUserInfo | undefined;
  setUserInfo(value?: ChatUserInfo): CreateChatRoomReq;
  hasUserInfo(): boolean;
  clearUserInfo(): CreateChatRoomReq;

  getRoomKey(): string;
  setRoomKey(value: string): CreateChatRoomReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateChatRoomReq.AsObject;
  static toObject(includeInstance: boolean, msg: CreateChatRoomReq): CreateChatRoomReq.AsObject;
  static serializeBinaryToWriter(message: CreateChatRoomReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateChatRoomReq;
  static deserializeBinaryFromReader(message: CreateChatRoomReq, reader: jspb.BinaryReader): CreateChatRoomReq;
}

export namespace CreateChatRoomReq {
  export type AsObject = {
    userInfo?: ChatUserInfo.AsObject,
    roomKey: string,
  }
}

export class CreateChatRoomRes extends jspb.Message {
  getError(): ChatError | undefined;
  setError(value?: ChatError): CreateChatRoomRes;
  hasError(): boolean;
  clearError(): CreateChatRoomRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CreateChatRoomRes.AsObject;
  static toObject(includeInstance: boolean, msg: CreateChatRoomRes): CreateChatRoomRes.AsObject;
  static serializeBinaryToWriter(message: CreateChatRoomRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CreateChatRoomRes;
  static deserializeBinaryFromReader(message: CreateChatRoomRes, reader: jspb.BinaryReader): CreateChatRoomRes;
}

export namespace CreateChatRoomRes {
  export type AsObject = {
    error?: ChatError.AsObject,
  }
}

export class QueryChatRoomInfoReq extends jspb.Message {
  getUserInfo(): ChatUserInfo | undefined;
  setUserInfo(value?: ChatUserInfo): QueryChatRoomInfoReq;
  hasUserInfo(): boolean;
  clearUserInfo(): QueryChatRoomInfoReq;

  getRoomKey(): string;
  setRoomKey(value: string): QueryChatRoomInfoReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryChatRoomInfoReq.AsObject;
  static toObject(includeInstance: boolean, msg: QueryChatRoomInfoReq): QueryChatRoomInfoReq.AsObject;
  static serializeBinaryToWriter(message: QueryChatRoomInfoReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryChatRoomInfoReq;
  static deserializeBinaryFromReader(message: QueryChatRoomInfoReq, reader: jspb.BinaryReader): QueryChatRoomInfoReq;
}

export namespace QueryChatRoomInfoReq {
  export type AsObject = {
    userInfo?: ChatUserInfo.AsObject,
    roomKey: string,
  }
}

export class QueryChatRoomInfoRes extends jspb.Message {
  getError(): ChatError | undefined;
  setError(value?: ChatError): QueryChatRoomInfoRes;
  hasError(): boolean;
  clearError(): QueryChatRoomInfoRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryChatRoomInfoRes.AsObject;
  static toObject(includeInstance: boolean, msg: QueryChatRoomInfoRes): QueryChatRoomInfoRes.AsObject;
  static serializeBinaryToWriter(message: QueryChatRoomInfoRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryChatRoomInfoRes;
  static deserializeBinaryFromReader(message: QueryChatRoomInfoRes, reader: jspb.BinaryReader): QueryChatRoomInfoRes;
}

export namespace QueryChatRoomInfoRes {
  export type AsObject = {
    error?: ChatError.AsObject,
  }
}

export class JoinRoomReq extends jspb.Message {
  getUserInfo(): ChatUserInfo | undefined;
  setUserInfo(value?: ChatUserInfo): JoinRoomReq;
  hasUserInfo(): boolean;
  clearUserInfo(): JoinRoomReq;

  getRoomKey(): string;
  setRoomKey(value: string): JoinRoomReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinRoomReq.AsObject;
  static toObject(includeInstance: boolean, msg: JoinRoomReq): JoinRoomReq.AsObject;
  static serializeBinaryToWriter(message: JoinRoomReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinRoomReq;
  static deserializeBinaryFromReader(message: JoinRoomReq, reader: jspb.BinaryReader): JoinRoomReq;
}

export namespace JoinRoomReq {
  export type AsObject = {
    userInfo?: ChatUserInfo.AsObject,
    roomKey: string,
  }
}

export class JoinRoomRes extends jspb.Message {
  getError(): ChatError | undefined;
  setError(value?: ChatError): JoinRoomRes;
  hasError(): boolean;
  clearError(): JoinRoomRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinRoomRes.AsObject;
  static toObject(includeInstance: boolean, msg: JoinRoomRes): JoinRoomRes.AsObject;
  static serializeBinaryToWriter(message: JoinRoomRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinRoomRes;
  static deserializeBinaryFromReader(message: JoinRoomRes, reader: jspb.BinaryReader): JoinRoomRes;
}

export namespace JoinRoomRes {
  export type AsObject = {
    error?: ChatError.AsObject,
  }
}

export class LeaveRoomReq extends jspb.Message {
  getUserInfo(): ChatUserInfo | undefined;
  setUserInfo(value?: ChatUserInfo): LeaveRoomReq;
  hasUserInfo(): boolean;
  clearUserInfo(): LeaveRoomReq;

  getRoomKey(): string;
  setRoomKey(value: string): LeaveRoomReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LeaveRoomReq.AsObject;
  static toObject(includeInstance: boolean, msg: LeaveRoomReq): LeaveRoomReq.AsObject;
  static serializeBinaryToWriter(message: LeaveRoomReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LeaveRoomReq;
  static deserializeBinaryFromReader(message: LeaveRoomReq, reader: jspb.BinaryReader): LeaveRoomReq;
}

export namespace LeaveRoomReq {
  export type AsObject = {
    userInfo?: ChatUserInfo.AsObject,
    roomKey: string,
  }
}

export class LeaveRoomRes extends jspb.Message {
  getError(): ChatError | undefined;
  setError(value?: ChatError): LeaveRoomRes;
  hasError(): boolean;
  clearError(): LeaveRoomRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LeaveRoomRes.AsObject;
  static toObject(includeInstance: boolean, msg: LeaveRoomRes): LeaveRoomRes.AsObject;
  static serializeBinaryToWriter(message: LeaveRoomRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LeaveRoomRes;
  static deserializeBinaryFromReader(message: LeaveRoomRes, reader: jspb.BinaryReader): LeaveRoomRes;
}

export namespace LeaveRoomRes {
  export type AsObject = {
    error?: ChatError.AsObject,
  }
}

export class Send2ChatRoomReq extends jspb.Message {
  getUserInfo(): ChatUserInfo | undefined;
  setUserInfo(value?: ChatUserInfo): Send2ChatRoomReq;
  hasUserInfo(): boolean;
  clearUserInfo(): Send2ChatRoomReq;

  getRoomKey(): string;
  setRoomKey(value: string): Send2ChatRoomReq;

  getMsg(): string;
  setMsg(value: string): Send2ChatRoomReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Send2ChatRoomReq.AsObject;
  static toObject(includeInstance: boolean, msg: Send2ChatRoomReq): Send2ChatRoomReq.AsObject;
  static serializeBinaryToWriter(message: Send2ChatRoomReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Send2ChatRoomReq;
  static deserializeBinaryFromReader(message: Send2ChatRoomReq, reader: jspb.BinaryReader): Send2ChatRoomReq;
}

export namespace Send2ChatRoomReq {
  export type AsObject = {
    userInfo?: ChatUserInfo.AsObject,
    roomKey: string,
    msg: string,
  }
}

export class Send2ChatRoomRes extends jspb.Message {
  getUserInfo(): ChatUserInfo | undefined;
  setUserInfo(value?: ChatUserInfo): Send2ChatRoomRes;
  hasUserInfo(): boolean;
  clearUserInfo(): Send2ChatRoomRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Send2ChatRoomRes.AsObject;
  static toObject(includeInstance: boolean, msg: Send2ChatRoomRes): Send2ChatRoomRes.AsObject;
  static serializeBinaryToWriter(message: Send2ChatRoomRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Send2ChatRoomRes;
  static deserializeBinaryFromReader(message: Send2ChatRoomRes, reader: jspb.BinaryReader): Send2ChatRoomRes;
}

export namespace Send2ChatRoomRes {
  export type AsObject = {
    userInfo?: ChatUserInfo.AsObject,
  }
}

