import * as jspb from 'google-protobuf'

import * as category_pb from './category_pb';
import * as video_pb from './video_pb';
import * as bt_pb from './bt_pb';
import * as google_api_annotations_pb from './google/api/annotations_pb';


export class UserInfo extends jspb.Message {
  getId(): number;
  setId(value: number): UserInfo;

  getName(): string;
  setName(value: string): UserInfo;

  getEmail(): string;
  setEmail(value: string): UserInfo;

  getPasswd(): string;
  setPasswd(value: string): UserInfo;

  getHomeDirectoryId(): number;
  setHomeDirectoryId(value: number): UserInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserInfo.AsObject;
  static toObject(includeInstance: boolean, msg: UserInfo): UserInfo.AsObject;
  static serializeBinaryToWriter(message: UserInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserInfo;
  static deserializeBinaryFromReader(message: UserInfo, reader: jspb.BinaryReader): UserInfo;
}

export namespace UserInfo {
  export type AsObject = {
    id: number,
    name: string,
    email: string,
    passwd: string,
    homeDirectoryId: number,
  }
}

export class RegisterInfo extends jspb.Message {
  getUserInfo(): UserInfo | undefined;
  setUserInfo(value?: UserInfo): RegisterInfo;
  hasUserInfo(): boolean;
  clearUserInfo(): RegisterInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterInfo.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterInfo): RegisterInfo.AsObject;
  static serializeBinaryToWriter(message: RegisterInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterInfo;
  static deserializeBinaryFromReader(message: RegisterInfo, reader: jspb.BinaryReader): RegisterInfo;
}

export namespace RegisterInfo {
  export type AsObject = {
    userInfo?: UserInfo.AsObject,
  }
}

export class RegisterRet extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterRet.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterRet): RegisterRet.AsObject;
  static serializeBinaryToWriter(message: RegisterRet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterRet;
  static deserializeBinaryFromReader(message: RegisterRet, reader: jspb.BinaryReader): RegisterRet;
}

export namespace RegisterRet {
  export type AsObject = {
  }
}

export class EmailInfo extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): EmailInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EmailInfo.AsObject;
  static toObject(includeInstance: boolean, msg: EmailInfo): EmailInfo.AsObject;
  static serializeBinaryToWriter(message: EmailInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EmailInfo;
  static deserializeBinaryFromReader(message: EmailInfo, reader: jspb.BinaryReader): EmailInfo;
}

export namespace EmailInfo {
  export type AsObject = {
    email: string,
  }
}

export class IsUsedEmailRet extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): IsUsedEmailRet.AsObject;
  static toObject(includeInstance: boolean, msg: IsUsedEmailRet): IsUsedEmailRet.AsObject;
  static serializeBinaryToWriter(message: IsUsedEmailRet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): IsUsedEmailRet;
  static deserializeBinaryFromReader(message: IsUsedEmailRet, reader: jspb.BinaryReader): IsUsedEmailRet;
}

export namespace IsUsedEmailRet {
  export type AsObject = {
  }
}

export class LoginInfo extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): LoginInfo;

  getPasswd(): string;
  setPasswd(value: string): LoginInfo;

  getRememberMe(): boolean;
  setRememberMe(value: boolean): LoginInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginInfo.AsObject;
  static toObject(includeInstance: boolean, msg: LoginInfo): LoginInfo.AsObject;
  static serializeBinaryToWriter(message: LoginInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginInfo;
  static deserializeBinaryFromReader(message: LoginInfo, reader: jspb.BinaryReader): LoginInfo;
}

export namespace LoginInfo {
  export type AsObject = {
    email: string,
    passwd: string,
    rememberMe: boolean,
  }
}

export class LoginRet extends jspb.Message {
  getToken(): string;
  setToken(value: string): LoginRet;

  getUserInfo(): UserInfo | undefined;
  setUserInfo(value?: UserInfo): LoginRet;
  hasUserInfo(): boolean;
  clearUserInfo(): LoginRet;

  getRememberMe(): boolean;
  setRememberMe(value: boolean): LoginRet;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginRet.AsObject;
  static toObject(includeInstance: boolean, msg: LoginRet): LoginRet.AsObject;
  static serializeBinaryToWriter(message: LoginRet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginRet;
  static deserializeBinaryFromReader(message: LoginRet, reader: jspb.BinaryReader): LoginRet;
}

export namespace LoginRet {
  export type AsObject = {
    token: string,
    userInfo?: UserInfo.AsObject,
    rememberMe: boolean,
  }
}

export class NewCategoryItemReq extends jspb.Message {
  getName(): string;
  setName(value: string): NewCategoryItemReq;

  getTypeId(): category_pb.CategoryItem.Type;
  setTypeId(value: category_pb.CategoryItem.Type): NewCategoryItemReq;

  getResourcePath(): string;
  setResourcePath(value: string): NewCategoryItemReq;

  getIntroduce(): string;
  setIntroduce(value: string): NewCategoryItemReq;

  getParentId(): number;
  setParentId(value: number): NewCategoryItemReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewCategoryItemReq.AsObject;
  static toObject(includeInstance: boolean, msg: NewCategoryItemReq): NewCategoryItemReq.AsObject;
  static serializeBinaryToWriter(message: NewCategoryItemReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewCategoryItemReq;
  static deserializeBinaryFromReader(message: NewCategoryItemReq, reader: jspb.BinaryReader): NewCategoryItemReq;
}

export namespace NewCategoryItemReq {
  export type AsObject = {
    name: string,
    typeId: category_pb.CategoryItem.Type,
    resourcePath: string,
    introduce: string,
    parentId: number,
  }
}

export class NewCategoryItemRes extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewCategoryItemRes.AsObject;
  static toObject(includeInstance: boolean, msg: NewCategoryItemRes): NewCategoryItemRes.AsObject;
  static serializeBinaryToWriter(message: NewCategoryItemRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewCategoryItemRes;
  static deserializeBinaryFromReader(message: NewCategoryItemRes, reader: jspb.BinaryReader): NewCategoryItemRes;
}

export namespace NewCategoryItemRes {
  export type AsObject = {
  }
}

export class DelCategoryItemReq extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): DelCategoryItemReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DelCategoryItemReq.AsObject;
  static toObject(includeInstance: boolean, msg: DelCategoryItemReq): DelCategoryItemReq.AsObject;
  static serializeBinaryToWriter(message: DelCategoryItemReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DelCategoryItemReq;
  static deserializeBinaryFromReader(message: DelCategoryItemReq, reader: jspb.BinaryReader): DelCategoryItemReq;
}

export namespace DelCategoryItemReq {
  export type AsObject = {
    itemId: number,
  }
}

export class DelCategoryItemRes extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DelCategoryItemRes.AsObject;
  static toObject(includeInstance: boolean, msg: DelCategoryItemRes): DelCategoryItemRes.AsObject;
  static serializeBinaryToWriter(message: DelCategoryItemRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DelCategoryItemRes;
  static deserializeBinaryFromReader(message: DelCategoryItemRes, reader: jspb.BinaryReader): DelCategoryItemRes;
}

export namespace DelCategoryItemRes {
  export type AsObject = {
  }
}

export class QuerySubItemsReq extends jspb.Message {
  getParentId(): number;
  setParentId(value: number): QuerySubItemsReq;

  getShareId(): string;
  setShareId(value: string): QuerySubItemsReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QuerySubItemsReq.AsObject;
  static toObject(includeInstance: boolean, msg: QuerySubItemsReq): QuerySubItemsReq.AsObject;
  static serializeBinaryToWriter(message: QuerySubItemsReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QuerySubItemsReq;
  static deserializeBinaryFromReader(message: QuerySubItemsReq, reader: jspb.BinaryReader): QuerySubItemsReq;
}

export namespace QuerySubItemsReq {
  export type AsObject = {
    parentId: number,
    shareId: string,
  }
}

export class QuerySubItemsRes extends jspb.Message {
  getParentItem(): category_pb.CategoryItem | undefined;
  setParentItem(value?: category_pb.CategoryItem): QuerySubItemsRes;
  hasParentItem(): boolean;
  clearParentItem(): QuerySubItemsRes;

  getItemsList(): Array<category_pb.CategoryItem>;
  setItemsList(value: Array<category_pb.CategoryItem>): QuerySubItemsRes;
  clearItemsList(): QuerySubItemsRes;
  addItems(value?: category_pb.CategoryItem, index?: number): category_pb.CategoryItem;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QuerySubItemsRes.AsObject;
  static toObject(includeInstance: boolean, msg: QuerySubItemsRes): QuerySubItemsRes.AsObject;
  static serializeBinaryToWriter(message: QuerySubItemsRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QuerySubItemsRes;
  static deserializeBinaryFromReader(message: QuerySubItemsRes, reader: jspb.BinaryReader): QuerySubItemsRes;
}

export namespace QuerySubItemsRes {
  export type AsObject = {
    parentItem?: category_pb.CategoryItem.AsObject,
    itemsList: Array<category_pb.CategoryItem.AsObject>,
  }
}

export class QueryBtVideosReq extends jspb.Message {
  getInfoHash(): bt_pb.InfoHash | undefined;
  setInfoHash(value?: bt_pb.InfoHash): QueryBtVideosReq;
  hasInfoHash(): boolean;
  clearInfoHash(): QueryBtVideosReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryBtVideosReq.AsObject;
  static toObject(includeInstance: boolean, msg: QueryBtVideosReq): QueryBtVideosReq.AsObject;
  static serializeBinaryToWriter(message: QueryBtVideosReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryBtVideosReq;
  static deserializeBinaryFromReader(message: QueryBtVideosReq, reader: jspb.BinaryReader): QueryBtVideosReq;
}

export namespace QueryBtVideosReq {
  export type AsObject = {
    infoHash?: bt_pb.InfoHash.AsObject,
  }
}

export class BtFileMetadata extends jspb.Message {
  getFileIndex(): number;
  setFileIndex(value: number): BtFileMetadata;

  getMeta(): video_pb.VideoMetadata | undefined;
  setMeta(value?: video_pb.VideoMetadata): BtFileMetadata;
  hasMeta(): boolean;
  clearMeta(): BtFileMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BtFileMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: BtFileMetadata): BtFileMetadata.AsObject;
  static serializeBinaryToWriter(message: BtFileMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BtFileMetadata;
  static deserializeBinaryFromReader(message: BtFileMetadata, reader: jspb.BinaryReader): BtFileMetadata;
}

export namespace BtFileMetadata {
  export type AsObject = {
    fileIndex: number,
    meta?: video_pb.VideoMetadata.AsObject,
  }
}

export class QueryBtVideosRes extends jspb.Message {
  getDataList(): Array<BtFileMetadata>;
  setDataList(value: Array<BtFileMetadata>): QueryBtVideosRes;
  clearDataList(): QueryBtVideosRes;
  addData(value?: BtFileMetadata, index?: number): BtFileMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryBtVideosRes.AsObject;
  static toObject(includeInstance: boolean, msg: QueryBtVideosRes): QueryBtVideosRes.AsObject;
  static serializeBinaryToWriter(message: QueryBtVideosRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryBtVideosRes;
  static deserializeBinaryFromReader(message: QueryBtVideosRes, reader: jspb.BinaryReader): QueryBtVideosRes;
}

export namespace QueryBtVideosRes {
  export type AsObject = {
    dataList: Array<BtFileMetadata.AsObject>,
  }
}

export class AddBtVideosReq extends jspb.Message {
  getInfoHash(): bt_pb.InfoHash | undefined;
  setInfoHash(value?: bt_pb.InfoHash): AddBtVideosReq;
  hasInfoHash(): boolean;
  clearInfoHash(): AddBtVideosReq;

  getFileIndexesList(): Array<number>;
  setFileIndexesList(value: Array<number>): AddBtVideosReq;
  clearFileIndexesList(): AddBtVideosReq;
  addFileIndexes(value: number, index?: number): AddBtVideosReq;

  getCategoryItemId(): number;
  setCategoryItemId(value: number): AddBtVideosReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddBtVideosReq.AsObject;
  static toObject(includeInstance: boolean, msg: AddBtVideosReq): AddBtVideosReq.AsObject;
  static serializeBinaryToWriter(message: AddBtVideosReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddBtVideosReq;
  static deserializeBinaryFromReader(message: AddBtVideosReq, reader: jspb.BinaryReader): AddBtVideosReq;
}

export namespace AddBtVideosReq {
  export type AsObject = {
    infoHash?: bt_pb.InfoHash.AsObject,
    fileIndexesList: Array<number>,
    categoryItemId: number,
  }
}

export class AddBtVideosRes extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddBtVideosRes.AsObject;
  static toObject(includeInstance: boolean, msg: AddBtVideosRes): AddBtVideosRes.AsObject;
  static serializeBinaryToWriter(message: AddBtVideosRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddBtVideosRes;
  static deserializeBinaryFromReader(message: AddBtVideosRes, reader: jspb.BinaryReader): AddBtVideosRes;
}

export namespace AddBtVideosRes {
  export type AsObject = {
  }
}

export class QueryItemInfoReq extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): QueryItemInfoReq;

  getShareId(): string;
  setShareId(value: string): QueryItemInfoReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryItemInfoReq.AsObject;
  static toObject(includeInstance: boolean, msg: QueryItemInfoReq): QueryItemInfoReq.AsObject;
  static serializeBinaryToWriter(message: QueryItemInfoReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryItemInfoReq;
  static deserializeBinaryFromReader(message: QueryItemInfoReq, reader: jspb.BinaryReader): QueryItemInfoReq;
}

export namespace QueryItemInfoReq {
  export type AsObject = {
    itemId: number,
    shareId: string,
  }
}

export class QueryItemInfoRes extends jspb.Message {
  getItemInfo(): category_pb.CategoryItem | undefined;
  setItemInfo(value?: category_pb.CategoryItem): QueryItemInfoRes;
  hasItemInfo(): boolean;
  clearItemInfo(): QueryItemInfoRes;

  getVideoInfo(): video_pb.Video | undefined;
  setVideoInfo(value?: video_pb.Video): QueryItemInfoRes;
  hasVideoInfo(): boolean;
  clearVideoInfo(): QueryItemInfoRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QueryItemInfoRes.AsObject;
  static toObject(includeInstance: boolean, msg: QueryItemInfoRes): QueryItemInfoRes.AsObject;
  static serializeBinaryToWriter(message: QueryItemInfoRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QueryItemInfoRes;
  static deserializeBinaryFromReader(message: QueryItemInfoRes, reader: jspb.BinaryReader): QueryItemInfoRes;
}

export namespace QueryItemInfoRes {
  export type AsObject = {
    itemInfo?: category_pb.CategoryItem.AsObject,
    videoInfo?: video_pb.Video.AsObject,
  }
}

export class ShareItemReq extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): ShareItemReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ShareItemReq.AsObject;
  static toObject(includeInstance: boolean, msg: ShareItemReq): ShareItemReq.AsObject;
  static serializeBinaryToWriter(message: ShareItemReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ShareItemReq;
  static deserializeBinaryFromReader(message: ShareItemReq, reader: jspb.BinaryReader): ShareItemReq;
}

export namespace ShareItemReq {
  export type AsObject = {
    itemId: number,
  }
}

export class ShareItemRes extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): ShareItemRes;

  getShareId(): string;
  setShareId(value: string): ShareItemRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ShareItemRes.AsObject;
  static toObject(includeInstance: boolean, msg: ShareItemRes): ShareItemRes.AsObject;
  static serializeBinaryToWriter(message: ShareItemRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ShareItemRes;
  static deserializeBinaryFromReader(message: ShareItemRes, reader: jspb.BinaryReader): ShareItemRes;
}

export namespace ShareItemRes {
  export type AsObject = {
    itemId: number,
    shareId: string,
  }
}

export class QuerySharedItemsReq extends jspb.Message {
  getUserId(): number;
  setUserId(value: number): QuerySharedItemsReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QuerySharedItemsReq.AsObject;
  static toObject(includeInstance: boolean, msg: QuerySharedItemsReq): QuerySharedItemsReq.AsObject;
  static serializeBinaryToWriter(message: QuerySharedItemsReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QuerySharedItemsReq;
  static deserializeBinaryFromReader(message: QuerySharedItemsReq, reader: jspb.BinaryReader): QuerySharedItemsReq;
}

export namespace QuerySharedItemsReq {
  export type AsObject = {
    userId: number,
  }
}

export class QuerySharedItemsRes extends jspb.Message {
  getSharedItemsList(): Array<category_pb.SharedItem>;
  setSharedItemsList(value: Array<category_pb.SharedItem>): QuerySharedItemsRes;
  clearSharedItemsList(): QuerySharedItemsRes;
  addSharedItems(value?: category_pb.SharedItem, index?: number): category_pb.SharedItem;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): QuerySharedItemsRes.AsObject;
  static toObject(includeInstance: boolean, msg: QuerySharedItemsRes): QuerySharedItemsRes.AsObject;
  static serializeBinaryToWriter(message: QuerySharedItemsRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): QuerySharedItemsRes;
  static deserializeBinaryFromReader(message: QuerySharedItemsRes, reader: jspb.BinaryReader): QuerySharedItemsRes;
}

export namespace QuerySharedItemsRes {
  export type AsObject = {
    sharedItemsList: Array<category_pb.SharedItem.AsObject>,
  }
}

export class DelSharedItemReq extends jspb.Message {
  getShareId(): string;
  setShareId(value: string): DelSharedItemReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DelSharedItemReq.AsObject;
  static toObject(includeInstance: boolean, msg: DelSharedItemReq): DelSharedItemReq.AsObject;
  static serializeBinaryToWriter(message: DelSharedItemReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DelSharedItemReq;
  static deserializeBinaryFromReader(message: DelSharedItemReq, reader: jspb.BinaryReader): DelSharedItemReq;
}

export namespace DelSharedItemReq {
  export type AsObject = {
    shareId: string,
  }
}

export class DelSharedItemRes extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DelSharedItemRes.AsObject;
  static toObject(includeInstance: boolean, msg: DelSharedItemRes): DelSharedItemRes.AsObject;
  static serializeBinaryToWriter(message: DelSharedItemRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DelSharedItemRes;
  static deserializeBinaryFromReader(message: DelSharedItemRes, reader: jspb.BinaryReader): DelSharedItemRes;
}

export namespace DelSharedItemRes {
  export type AsObject = {
  }
}

export class RefreshSubtitleReq extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): RefreshSubtitleReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshSubtitleReq.AsObject;
  static toObject(includeInstance: boolean, msg: RefreshSubtitleReq): RefreshSubtitleReq.AsObject;
  static serializeBinaryToWriter(message: RefreshSubtitleReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RefreshSubtitleReq;
  static deserializeBinaryFromReader(message: RefreshSubtitleReq, reader: jspb.BinaryReader): RefreshSubtitleReq;
}

export namespace RefreshSubtitleReq {
  export type AsObject = {
    itemId: number,
  }
}

export class RefreshSubtitleRes extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshSubtitleRes.AsObject;
  static toObject(includeInstance: boolean, msg: RefreshSubtitleRes): RefreshSubtitleRes.AsObject;
  static serializeBinaryToWriter(message: RefreshSubtitleRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RefreshSubtitleRes;
  static deserializeBinaryFromReader(message: RefreshSubtitleRes, reader: jspb.BinaryReader): RefreshSubtitleRes;
}

export namespace RefreshSubtitleRes {
  export type AsObject = {
  }
}

export class SubtitleFile extends jspb.Message {
  getName(): string;
  setName(value: string): SubtitleFile;

  getContent(): Uint8Array | string;
  getContent_asU8(): Uint8Array;
  getContent_asB64(): string;
  setContent(value: Uint8Array | string): SubtitleFile;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SubtitleFile.AsObject;
  static toObject(includeInstance: boolean, msg: SubtitleFile): SubtitleFile.AsObject;
  static serializeBinaryToWriter(message: SubtitleFile, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SubtitleFile;
  static deserializeBinaryFromReader(message: SubtitleFile, reader: jspb.BinaryReader): SubtitleFile;
}

export namespace SubtitleFile {
  export type AsObject = {
    name: string,
    content: Uint8Array | string,
  }
}

export class UploadSubtitleReq extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): UploadSubtitleReq;

  getSubtitlesList(): Array<SubtitleFile>;
  setSubtitlesList(value: Array<SubtitleFile>): UploadSubtitleReq;
  clearSubtitlesList(): UploadSubtitleReq;
  addSubtitles(value?: SubtitleFile, index?: number): SubtitleFile;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadSubtitleReq.AsObject;
  static toObject(includeInstance: boolean, msg: UploadSubtitleReq): UploadSubtitleReq.AsObject;
  static serializeBinaryToWriter(message: UploadSubtitleReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadSubtitleReq;
  static deserializeBinaryFromReader(message: UploadSubtitleReq, reader: jspb.BinaryReader): UploadSubtitleReq;
}

export namespace UploadSubtitleReq {
  export type AsObject = {
    itemId: number,
    subtitlesList: Array<SubtitleFile.AsObject>,
  }
}

export class UploadSubtitleRes extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadSubtitleRes.AsObject;
  static toObject(includeInstance: boolean, msg: UploadSubtitleRes): UploadSubtitleRes.AsObject;
  static serializeBinaryToWriter(message: UploadSubtitleRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadSubtitleRes;
  static deserializeBinaryFromReader(message: UploadSubtitleRes, reader: jspb.BinaryReader): UploadSubtitleRes;
}

export namespace UploadSubtitleRes {
  export type AsObject = {
  }
}

export class JoinChatRoomReq extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): JoinChatRoomReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinChatRoomReq.AsObject;
  static toObject(includeInstance: boolean, msg: JoinChatRoomReq): JoinChatRoomReq.AsObject;
  static serializeBinaryToWriter(message: JoinChatRoomReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinChatRoomReq;
  static deserializeBinaryFromReader(message: JoinChatRoomReq, reader: jspb.BinaryReader): JoinChatRoomReq;
}

export namespace JoinChatRoomReq {
  export type AsObject = {
    itemId: number,
  }
}

export class ChatMessage extends jspb.Message {
  getUserId(): number;
  setUserId(value: number): ChatMessage;

  getUserName(): string;
  setUserName(value: string): ChatMessage;

  getSentTime(): number;
  setSentTime(value: number): ChatMessage;

  getMsg(): string;
  setMsg(value: string): ChatMessage;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChatMessage.AsObject;
  static toObject(includeInstance: boolean, msg: ChatMessage): ChatMessage.AsObject;
  static serializeBinaryToWriter(message: ChatMessage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChatMessage;
  static deserializeBinaryFromReader(message: ChatMessage, reader: jspb.BinaryReader): ChatMessage;
}

export namespace ChatMessage {
  export type AsObject = {
    userId: number,
    userName: string,
    sentTime: number,
    msg: string,
  }
}

export class JoinChatRoomRes extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): JoinChatRoomRes;

  getChatMsgsList(): Array<ChatMessage>;
  setChatMsgsList(value: Array<ChatMessage>): JoinChatRoomRes;
  clearChatMsgsList(): JoinChatRoomRes;
  addChatMsgs(value?: ChatMessage, index?: number): ChatMessage;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinChatRoomRes.AsObject;
  static toObject(includeInstance: boolean, msg: JoinChatRoomRes): JoinChatRoomRes.AsObject;
  static serializeBinaryToWriter(message: JoinChatRoomRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinChatRoomRes;
  static deserializeBinaryFromReader(message: JoinChatRoomRes, reader: jspb.BinaryReader): JoinChatRoomRes;
}

export namespace JoinChatRoomRes {
  export type AsObject = {
    itemId: number,
    chatMsgsList: Array<ChatMessage.AsObject>,
  }
}

export class SendMsg2ChatRoomReq extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): SendMsg2ChatRoomReq;

  getChatMsg(): ChatMessage | undefined;
  setChatMsg(value?: ChatMessage): SendMsg2ChatRoomReq;
  hasChatMsg(): boolean;
  clearChatMsg(): SendMsg2ChatRoomReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendMsg2ChatRoomReq.AsObject;
  static toObject(includeInstance: boolean, msg: SendMsg2ChatRoomReq): SendMsg2ChatRoomReq.AsObject;
  static serializeBinaryToWriter(message: SendMsg2ChatRoomReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendMsg2ChatRoomReq;
  static deserializeBinaryFromReader(message: SendMsg2ChatRoomReq, reader: jspb.BinaryReader): SendMsg2ChatRoomReq;
}

export namespace SendMsg2ChatRoomReq {
  export type AsObject = {
    itemId: number,
    chatMsg?: ChatMessage.AsObject,
  }
}

export class SendMsg2ChatRoomRes extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendMsg2ChatRoomRes.AsObject;
  static toObject(includeInstance: boolean, msg: SendMsg2ChatRoomRes): SendMsg2ChatRoomRes.AsObject;
  static serializeBinaryToWriter(message: SendMsg2ChatRoomRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendMsg2ChatRoomRes;
  static deserializeBinaryFromReader(message: SendMsg2ChatRoomRes, reader: jspb.BinaryReader): SendMsg2ChatRoomRes;
}

export namespace SendMsg2ChatRoomRes {
  export type AsObject = {
  }
}

