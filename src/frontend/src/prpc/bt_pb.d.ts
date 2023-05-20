import * as jspb from 'google-protobuf'



export class DownloadRequest extends jspb.Message {
  getType(): DownloadRequest.ReqType;
  setType(value: DownloadRequest.ReqType): DownloadRequest;

  getContent(): Uint8Array | string;
  getContent_asU8(): Uint8Array;
  getContent_asB64(): string;
  setContent(value: Uint8Array | string): DownloadRequest;

  getSavePath(): string;
  setSavePath(value: string): DownloadRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DownloadRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DownloadRequest): DownloadRequest.AsObject;
  static serializeBinaryToWriter(message: DownloadRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DownloadRequest;
  static deserializeBinaryFromReader(message: DownloadRequest, reader: jspb.BinaryReader): DownloadRequest;
}

export namespace DownloadRequest {
  export type AsObject = {
    type: DownloadRequest.ReqType,
    content: Uint8Array | string,
    savePath: string,
  }

  export enum ReqType { 
    MAGNETURI = 0,
    TORRENT = 1,
    RESUME = 2,
  }
}

export class InfoHash extends jspb.Message {
  getVersion(): number;
  setVersion(value: number): InfoHash;

  getHash(): Uint8Array | string;
  getHash_asU8(): Uint8Array;
  getHash_asB64(): string;
  setHash(value: Uint8Array | string): InfoHash;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InfoHash.AsObject;
  static toObject(includeInstance: boolean, msg: InfoHash): InfoHash.AsObject;
  static serializeBinaryToWriter(message: InfoHash, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InfoHash;
  static deserializeBinaryFromReader(message: InfoHash, reader: jspb.BinaryReader): InfoHash;
}

export namespace InfoHash {
  export type AsObject = {
    version: number,
    hash: Uint8Array | string,
  }
}

export class DownloadRespone extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): DownloadRespone;
  hasInfoHash(): boolean;
  clearInfoHash(): DownloadRespone;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DownloadRespone.AsObject;
  static toObject(includeInstance: boolean, msg: DownloadRespone): DownloadRespone.AsObject;
  static serializeBinaryToWriter(message: DownloadRespone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DownloadRespone;
  static deserializeBinaryFromReader(message: DownloadRespone, reader: jspb.BinaryReader): DownloadRespone;
}

export namespace DownloadRespone {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
  }
}

export class StatusRequest extends jspb.Message {
  getInfoHashList(): Array<InfoHash>;
  setInfoHashList(value: Array<InfoHash>): StatusRequest;
  clearInfoHashList(): StatusRequest;
  addInfoHash(value?: InfoHash, index?: number): InfoHash;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StatusRequest): StatusRequest.AsObject;
  static serializeBinaryToWriter(message: StatusRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusRequest;
  static deserializeBinaryFromReader(message: StatusRequest, reader: jspb.BinaryReader): StatusRequest;
}

export namespace StatusRequest {
  export type AsObject = {
    infoHashList: Array<InfoHash.AsObject>,
  }
}

export class TorrentStatus extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): TorrentStatus;
  hasInfoHash(): boolean;
  clearInfoHash(): TorrentStatus;

  getName(): string;
  setName(value: string): TorrentStatus;

  getDownloadPayloadRate(): number;
  setDownloadPayloadRate(value: number): TorrentStatus;

  getTotalDone(): number;
  setTotalDone(value: number): TorrentStatus;

  getTotal(): number;
  setTotal(value: number): TorrentStatus;

  getProgress(): number;
  setProgress(value: number): TorrentStatus;

  getNumPeers(): number;
  setNumPeers(value: number): TorrentStatus;

  getState(): BtStateEnum;
  setState(value: BtStateEnum): TorrentStatus;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TorrentStatus.AsObject;
  static toObject(includeInstance: boolean, msg: TorrentStatus): TorrentStatus.AsObject;
  static serializeBinaryToWriter(message: TorrentStatus, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TorrentStatus;
  static deserializeBinaryFromReader(message: TorrentStatus, reader: jspb.BinaryReader): TorrentStatus;
}

export namespace TorrentStatus {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
    name: string,
    downloadPayloadRate: number,
    totalDone: number,
    total: number,
    progress: number,
    numPeers: number,
    state: BtStateEnum,
  }
}

export class StatusRespone extends jspb.Message {
  getStatusArrayList(): Array<TorrentStatus>;
  setStatusArrayList(value: Array<TorrentStatus>): StatusRespone;
  clearStatusArrayList(): StatusRespone;
  addStatusArray(value?: TorrentStatus, index?: number): TorrentStatus;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StatusRespone.AsObject;
  static toObject(includeInstance: boolean, msg: StatusRespone): StatusRespone.AsObject;
  static serializeBinaryToWriter(message: StatusRespone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StatusRespone;
  static deserializeBinaryFromReader(message: StatusRespone, reader: jspb.BinaryReader): StatusRespone;
}

export namespace StatusRespone {
  export type AsObject = {
    statusArrayList: Array<TorrentStatus.AsObject>,
  }
}

export class TorrentInfoReq extends jspb.Message {
  getInfoHashList(): Array<InfoHash>;
  setInfoHashList(value: Array<InfoHash>): TorrentInfoReq;
  clearInfoHashList(): TorrentInfoReq;
  addInfoHash(value?: InfoHash, index?: number): InfoHash;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TorrentInfoReq.AsObject;
  static toObject(includeInstance: boolean, msg: TorrentInfoReq): TorrentInfoReq.AsObject;
  static serializeBinaryToWriter(message: TorrentInfoReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TorrentInfoReq;
  static deserializeBinaryFromReader(message: TorrentInfoReq, reader: jspb.BinaryReader): TorrentInfoReq;
}

export namespace TorrentInfoReq {
  export type AsObject = {
    infoHashList: Array<InfoHash.AsObject>,
  }
}

export class BtFile extends jspb.Message {
  getName(): string;
  setName(value: string): BtFile;

  getIndex(): number;
  setIndex(value: number): BtFile;

  getSt(): BtFile.State;
  setSt(value: BtFile.State): BtFile;

  getTotalSize(): number;
  setTotalSize(value: number): BtFile;

  getDownloaded(): number;
  setDownloaded(value: number): BtFile;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BtFile.AsObject;
  static toObject(includeInstance: boolean, msg: BtFile): BtFile.AsObject;
  static serializeBinaryToWriter(message: BtFile, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BtFile;
  static deserializeBinaryFromReader(message: BtFile, reader: jspb.BinaryReader): BtFile;
}

export namespace BtFile {
  export type AsObject = {
    name: string,
    index: number,
    st: BtFile.State,
    totalSize: number,
    downloaded: number,
  }

  export enum State { 
    STOP = 0,
    DOWNLOADING = 1,
    COMPLETED = 2,
  }
}

export class TorrentInfo extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): TorrentInfo;
  hasInfoHash(): boolean;
  clearInfoHash(): TorrentInfo;

  getName(): string;
  setName(value: string): TorrentInfo;

  getState(): BtStateEnum;
  setState(value: BtStateEnum): TorrentInfo;

  getSavePath(): string;
  setSavePath(value: string): TorrentInfo;

  getFilesList(): Array<BtFile>;
  setFilesList(value: Array<BtFile>): TorrentInfo;
  clearFilesList(): TorrentInfo;
  addFiles(value?: BtFile, index?: number): BtFile;

  getTotalSize(): number;
  setTotalSize(value: number): TorrentInfo;

  getPieceLength(): number;
  setPieceLength(value: number): TorrentInfo;

  getNumPieces(): number;
  setNumPieces(value: number): TorrentInfo;

  getResumeData(): Uint8Array | string;
  getResumeData_asU8(): Uint8Array;
  getResumeData_asB64(): string;
  setResumeData(value: Uint8Array | string): TorrentInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TorrentInfo.AsObject;
  static toObject(includeInstance: boolean, msg: TorrentInfo): TorrentInfo.AsObject;
  static serializeBinaryToWriter(message: TorrentInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TorrentInfo;
  static deserializeBinaryFromReader(message: TorrentInfo, reader: jspb.BinaryReader): TorrentInfo;
}

export namespace TorrentInfo {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
    name: string,
    state: BtStateEnum,
    savePath: string,
    filesList: Array<BtFile.AsObject>,
    totalSize: number,
    pieceLength: number,
    numPieces: number,
    resumeData: Uint8Array | string,
  }
}

export class TorrentInfoRes extends jspb.Message {
  getTi(): TorrentInfo | undefined;
  setTi(value?: TorrentInfo): TorrentInfoRes;
  hasTi(): boolean;
  clearTi(): TorrentInfoRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TorrentInfoRes.AsObject;
  static toObject(includeInstance: boolean, msg: TorrentInfoRes): TorrentInfoRes.AsObject;
  static serializeBinaryToWriter(message: TorrentInfoRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TorrentInfoRes;
  static deserializeBinaryFromReader(message: TorrentInfoRes, reader: jspb.BinaryReader): TorrentInfoRes;
}

export namespace TorrentInfoRes {
  export type AsObject = {
    ti?: TorrentInfo.AsObject,
  }
}

export class RemoveTorrentReq extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): RemoveTorrentReq;
  hasInfoHash(): boolean;
  clearInfoHash(): RemoveTorrentReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveTorrentReq.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveTorrentReq): RemoveTorrentReq.AsObject;
  static serializeBinaryToWriter(message: RemoveTorrentReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveTorrentReq;
  static deserializeBinaryFromReader(message: RemoveTorrentReq, reader: jspb.BinaryReader): RemoveTorrentReq;
}

export namespace RemoveTorrentReq {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
  }
}

export class RemoveTorrentRes extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveTorrentRes.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveTorrentRes): RemoveTorrentRes.AsObject;
  static serializeBinaryToWriter(message: RemoveTorrentRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveTorrentRes;
  static deserializeBinaryFromReader(message: RemoveTorrentRes, reader: jspb.BinaryReader): RemoveTorrentRes;
}

export namespace RemoveTorrentRes {
  export type AsObject = {
  }
}

export class FileProgressReq extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): FileProgressReq;
  hasInfoHash(): boolean;
  clearInfoHash(): FileProgressReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileProgressReq.AsObject;
  static toObject(includeInstance: boolean, msg: FileProgressReq): FileProgressReq.AsObject;
  static serializeBinaryToWriter(message: FileProgressReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileProgressReq;
  static deserializeBinaryFromReader(message: FileProgressReq, reader: jspb.BinaryReader): FileProgressReq;
}

export namespace FileProgressReq {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
  }
}

export class FileProgressRes extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): FileProgressRes;
  hasInfoHash(): boolean;
  clearInfoHash(): FileProgressRes;

  getFilesList(): Array<BtFile>;
  setFilesList(value: Array<BtFile>): FileProgressRes;
  clearFilesList(): FileProgressRes;
  addFiles(value?: BtFile, index?: number): BtFile;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileProgressRes.AsObject;
  static toObject(includeInstance: boolean, msg: FileProgressRes): FileProgressRes.AsObject;
  static serializeBinaryToWriter(message: FileProgressRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileProgressRes;
  static deserializeBinaryFromReader(message: FileProgressRes, reader: jspb.BinaryReader): FileProgressRes;
}

export namespace FileProgressRes {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
    filesList: Array<BtFile.AsObject>,
  }
}

export class FileCompletedReq extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileCompletedReq.AsObject;
  static toObject(includeInstance: boolean, msg: FileCompletedReq): FileCompletedReq.AsObject;
  static serializeBinaryToWriter(message: FileCompletedReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileCompletedReq;
  static deserializeBinaryFromReader(message: FileCompletedReq, reader: jspb.BinaryReader): FileCompletedReq;
}

export namespace FileCompletedReq {
  export type AsObject = {
  }
}

export class FileCompletedRes extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): FileCompletedRes;
  hasInfoHash(): boolean;
  clearInfoHash(): FileCompletedRes;

  getFileIndex(): number;
  setFileIndex(value: number): FileCompletedRes;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileCompletedRes.AsObject;
  static toObject(includeInstance: boolean, msg: FileCompletedRes): FileCompletedRes.AsObject;
  static serializeBinaryToWriter(message: FileCompletedRes, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileCompletedRes;
  static deserializeBinaryFromReader(message: FileCompletedRes, reader: jspb.BinaryReader): FileCompletedRes;
}

export namespace FileCompletedRes {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
    fileIndex: number,
  }
}

export enum BtStateEnum { 
  UNKNOWN = 0,
  CHECKING_FILES = 1,
  DOWNLOADING_METADATA = 2,
  DOWNLOADING = 3,
  FINISHED = 4,
  SEEDING = 5,
  CHECKING_RESUME_DATA = 7,
}
