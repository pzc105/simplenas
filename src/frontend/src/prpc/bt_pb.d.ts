import * as jspb from 'google-protobuf'



export class InitedSessionReq extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitedSessionReq.AsObject;
  static toObject(includeInstance: boolean, msg: InitedSessionReq): InitedSessionReq.AsObject;
  static serializeBinaryToWriter(message: InitedSessionReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitedSessionReq;
  static deserializeBinaryFromReader(message: InitedSessionReq, reader: jspb.BinaryReader): InitedSessionReq;
}

export namespace InitedSessionReq {
  export type AsObject = {
  }
}

export class InitedSessionRsp extends jspb.Message {
  getInited(): boolean;
  setInited(value: boolean): InitedSessionRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitedSessionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: InitedSessionRsp): InitedSessionRsp.AsObject;
  static serializeBinaryToWriter(message: InitedSessionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitedSessionRsp;
  static deserializeBinaryFromReader(message: InitedSessionRsp, reader: jspb.BinaryReader): InitedSessionRsp;
}

export namespace InitedSessionRsp {
  export type AsObject = {
    inited: boolean,
  }
}

export class InitSessionReq extends jspb.Message {
  getProxyHost(): string;
  setProxyHost(value: string): InitSessionReq;

  getProxyPort(): number;
  setProxyPort(value: number): InitSessionReq;

  getProxyType(): string;
  setProxyType(value: string): InitSessionReq;

  getUploadRateLimit(): number;
  setUploadRateLimit(value: number): InitSessionReq;

  getDownloadRateLimit(): number;
  setDownloadRateLimit(value: number): InitSessionReq;

  getHashingThreads(): number;
  setHashingThreads(value: number): InitSessionReq;

  getResumeData(): Uint8Array | string;
  getResumeData_asU8(): Uint8Array;
  getResumeData_asB64(): string;
  setResumeData(value: Uint8Array | string): InitSessionReq;

  getListenInterfaces(): string;
  setListenInterfaces(value: string): InitSessionReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitSessionReq.AsObject;
  static toObject(includeInstance: boolean, msg: InitSessionReq): InitSessionReq.AsObject;
  static serializeBinaryToWriter(message: InitSessionReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitSessionReq;
  static deserializeBinaryFromReader(message: InitSessionReq, reader: jspb.BinaryReader): InitSessionReq;
}

export namespace InitSessionReq {
  export type AsObject = {
    proxyHost: string,
    proxyPort: number,
    proxyType: string,
    uploadRateLimit: number,
    downloadRateLimit: number,
    hashingThreads: number,
    resumeData: Uint8Array | string,
    listenInterfaces: string,
  }
}

export class InitSessionRsp extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InitSessionRsp.AsObject;
  static toObject(includeInstance: boolean, msg: InitSessionRsp): InitSessionRsp.AsObject;
  static serializeBinaryToWriter(message: InitSessionRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InitSessionRsp;
  static deserializeBinaryFromReader(message: InitSessionRsp, reader: jspb.BinaryReader): InitSessionRsp;
}

export namespace InitSessionRsp {
  export type AsObject = {
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
    savePath: string,
    filesList: Array<BtFile.AsObject>,
    totalSize: number,
    pieceLength: number,
    numPieces: number,
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

export class DownloadRequest extends jspb.Message {
  getType(): DownloadRequest.ReqType;
  setType(value: DownloadRequest.ReqType): DownloadRequest;

  getContent(): Uint8Array | string;
  getContent_asU8(): Uint8Array;
  getContent_asB64(): string;
  setContent(value: Uint8Array | string): DownloadRequest;

  getSavePath(): string;
  setSavePath(value: string): DownloadRequest;

  getStopAfterGotMeta(): boolean;
  setStopAfterGotMeta(value: boolean): DownloadRequest;

  getTrackersList(): Array<string>;
  setTrackersList(value: Array<string>): DownloadRequest;
  clearTrackersList(): DownloadRequest;
  addTrackers(value: string, index?: number): DownloadRequest;

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
    stopAfterGotMeta: boolean,
    trackersList: Array<string>,
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

export class GetMagnetUriReq extends jspb.Message {
  getType(): GetMagnetUriReq.ReqType;
  setType(value: GetMagnetUriReq.ReqType): GetMagnetUriReq;

  getContent(): Uint8Array | string;
  getContent_asU8(): Uint8Array;
  getContent_asB64(): string;
  setContent(value: Uint8Array | string): GetMagnetUriReq;

  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): GetMagnetUriReq;
  hasInfoHash(): boolean;
  clearInfoHash(): GetMagnetUriReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMagnetUriReq.AsObject;
  static toObject(includeInstance: boolean, msg: GetMagnetUriReq): GetMagnetUriReq.AsObject;
  static serializeBinaryToWriter(message: GetMagnetUriReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMagnetUriReq;
  static deserializeBinaryFromReader(message: GetMagnetUriReq, reader: jspb.BinaryReader): GetMagnetUriReq;
}

export namespace GetMagnetUriReq {
  export type AsObject = {
    type: GetMagnetUriReq.ReqType,
    content: Uint8Array | string,
    infoHash?: InfoHash.AsObject,
  }

  export enum ReqType { 
    TORRENT = 0,
    INFOHASH = 1,
  }
}

export class GetMagnetUriRsp extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): GetMagnetUriRsp;
  hasInfoHash(): boolean;
  clearInfoHash(): GetMagnetUriRsp;

  getMagnetUri(): string;
  setMagnetUri(value: string): GetMagnetUriRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMagnetUriRsp.AsObject;
  static toObject(includeInstance: boolean, msg: GetMagnetUriRsp): GetMagnetUriRsp.AsObject;
  static serializeBinaryToWriter(message: GetMagnetUriRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMagnetUriRsp;
  static deserializeBinaryFromReader(message: GetMagnetUriRsp, reader: jspb.BinaryReader): GetMagnetUriRsp;
}

export namespace GetMagnetUriRsp {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
    magnetUri: string,
  }
}

export class GetResumeDataReq extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): GetResumeDataReq;
  hasInfoHash(): boolean;
  clearInfoHash(): GetResumeDataReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetResumeDataReq.AsObject;
  static toObject(includeInstance: boolean, msg: GetResumeDataReq): GetResumeDataReq.AsObject;
  static serializeBinaryToWriter(message: GetResumeDataReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetResumeDataReq;
  static deserializeBinaryFromReader(message: GetResumeDataReq, reader: jspb.BinaryReader): GetResumeDataReq;
}

export namespace GetResumeDataReq {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
  }
}

export class GetResumeDataRsp extends jspb.Message {
  getResumeData(): Uint8Array | string;
  getResumeData_asU8(): Uint8Array;
  getResumeData_asB64(): string;
  setResumeData(value: Uint8Array | string): GetResumeDataRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetResumeDataRsp.AsObject;
  static toObject(includeInstance: boolean, msg: GetResumeDataRsp): GetResumeDataRsp.AsObject;
  static serializeBinaryToWriter(message: GetResumeDataRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetResumeDataRsp;
  static deserializeBinaryFromReader(message: GetResumeDataRsp, reader: jspb.BinaryReader): GetResumeDataRsp;
}

export namespace GetResumeDataRsp {
  export type AsObject = {
    resumeData: Uint8Array | string,
  }
}

export class GetTorrentInfoReq extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): GetTorrentInfoReq;
  hasInfoHash(): boolean;
  clearInfoHash(): GetTorrentInfoReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTorrentInfoReq.AsObject;
  static toObject(includeInstance: boolean, msg: GetTorrentInfoReq): GetTorrentInfoReq.AsObject;
  static serializeBinaryToWriter(message: GetTorrentInfoReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTorrentInfoReq;
  static deserializeBinaryFromReader(message: GetTorrentInfoReq, reader: jspb.BinaryReader): GetTorrentInfoReq;
}

export namespace GetTorrentInfoReq {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
  }
}

export class GetTorrentInfoRsp extends jspb.Message {
  getTorrentInfo(): TorrentInfo | undefined;
  setTorrentInfo(value?: TorrentInfo): GetTorrentInfoRsp;
  hasTorrentInfo(): boolean;
  clearTorrentInfo(): GetTorrentInfoRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetTorrentInfoRsp.AsObject;
  static toObject(includeInstance: boolean, msg: GetTorrentInfoRsp): GetTorrentInfoRsp.AsObject;
  static serializeBinaryToWriter(message: GetTorrentInfoRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetTorrentInfoRsp;
  static deserializeBinaryFromReader(message: GetTorrentInfoRsp, reader: jspb.BinaryReader): GetTorrentInfoRsp;
}

export namespace GetTorrentInfoRsp {
  export type AsObject = {
    torrentInfo?: TorrentInfo.AsObject,
  }
}

export class GetSessionParamsReq extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSessionParamsReq.AsObject;
  static toObject(includeInstance: boolean, msg: GetSessionParamsReq): GetSessionParamsReq.AsObject;
  static serializeBinaryToWriter(message: GetSessionParamsReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSessionParamsReq;
  static deserializeBinaryFromReader(message: GetSessionParamsReq, reader: jspb.BinaryReader): GetSessionParamsReq;
}

export namespace GetSessionParamsReq {
  export type AsObject = {
  }
}

export class GetSessionParamsRsp extends jspb.Message {
  getResumeData(): Uint8Array | string;
  getResumeData_asU8(): Uint8Array;
  getResumeData_asB64(): string;
  setResumeData(value: Uint8Array | string): GetSessionParamsRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetSessionParamsRsp.AsObject;
  static toObject(includeInstance: boolean, msg: GetSessionParamsRsp): GetSessionParamsRsp.AsObject;
  static serializeBinaryToWriter(message: GetSessionParamsRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetSessionParamsRsp;
  static deserializeBinaryFromReader(message: GetSessionParamsRsp, reader: jspb.BinaryReader): GetSessionParamsRsp;
}

export namespace GetSessionParamsRsp {
  export type AsObject = {
    resumeData: Uint8Array | string,
  }
}

export class GetBtStatusReq extends jspb.Message {
  getInfoHash(): InfoHash | undefined;
  setInfoHash(value?: InfoHash): GetBtStatusReq;
  hasInfoHash(): boolean;
  clearInfoHash(): GetBtStatusReq;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetBtStatusReq.AsObject;
  static toObject(includeInstance: boolean, msg: GetBtStatusReq): GetBtStatusReq.AsObject;
  static serializeBinaryToWriter(message: GetBtStatusReq, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetBtStatusReq;
  static deserializeBinaryFromReader(message: GetBtStatusReq, reader: jspb.BinaryReader): GetBtStatusReq;
}

export namespace GetBtStatusReq {
  export type AsObject = {
    infoHash?: InfoHash.AsObject,
  }
}

export class GetBtStatusRsp extends jspb.Message {
  getStatus(): TorrentStatus | undefined;
  setStatus(value?: TorrentStatus): GetBtStatusRsp;
  hasStatus(): boolean;
  clearStatus(): GetBtStatusRsp;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetBtStatusRsp.AsObject;
  static toObject(includeInstance: boolean, msg: GetBtStatusRsp): GetBtStatusRsp.AsObject;
  static serializeBinaryToWriter(message: GetBtStatusRsp, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetBtStatusRsp;
  static deserializeBinaryFromReader(message: GetBtStatusRsp, reader: jspb.BinaryReader): GetBtStatusRsp;
}

export namespace GetBtStatusRsp {
  export type AsObject = {
    status?: TorrentStatus.AsObject,
  }
}

export class BtStatusRequest extends jspb.Message {
  getInfoHashList(): Array<InfoHash>;
  setInfoHashList(value: Array<InfoHash>): BtStatusRequest;
  clearInfoHashList(): BtStatusRequest;
  addInfoHash(value?: InfoHash, index?: number): InfoHash;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BtStatusRequest.AsObject;
  static toObject(includeInstance: boolean, msg: BtStatusRequest): BtStatusRequest.AsObject;
  static serializeBinaryToWriter(message: BtStatusRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BtStatusRequest;
  static deserializeBinaryFromReader(message: BtStatusRequest, reader: jspb.BinaryReader): BtStatusRequest;
}

export namespace BtStatusRequest {
  export type AsObject = {
    infoHashList: Array<InfoHash.AsObject>,
  }
}

export class BtStatusRespone extends jspb.Message {
  getStatusArrayList(): Array<TorrentStatus>;
  setStatusArrayList(value: Array<TorrentStatus>): BtStatusRespone;
  clearStatusArrayList(): BtStatusRespone;
  addStatusArray(value?: TorrentStatus, index?: number): TorrentStatus;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BtStatusRespone.AsObject;
  static toObject(includeInstance: boolean, msg: BtStatusRespone): BtStatusRespone.AsObject;
  static serializeBinaryToWriter(message: BtStatusRespone, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BtStatusRespone;
  static deserializeBinaryFromReader(message: BtStatusRespone, reader: jspb.BinaryReader): BtStatusRespone;
}

export namespace BtStatusRespone {
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
