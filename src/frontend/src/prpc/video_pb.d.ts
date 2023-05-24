import * as jspb from 'google-protobuf'

import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class VideoStream extends jspb.Message {
  getIndex(): number;
  setIndex(value: number): VideoStream;

  getCodecName(): string;
  setCodecName(value: string): VideoStream;

  getCodecLongName(): string;
  setCodecLongName(value: string): VideoStream;

  getProfile(): string;
  setProfile(value: string): VideoStream;

  getCodecType(): string;
  setCodecType(value: string): VideoStream;

  getWidth(): number;
  setWidth(value: number): VideoStream;

  getHeight(): number;
  setHeight(value: number): VideoStream;

  getRFrameRate(): string;
  setRFrameRate(value: string): VideoStream;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VideoStream.AsObject;
  static toObject(includeInstance: boolean, msg: VideoStream): VideoStream.AsObject;
  static serializeBinaryToWriter(message: VideoStream, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VideoStream;
  static deserializeBinaryFromReader(message: VideoStream, reader: jspb.BinaryReader): VideoStream;
}

export namespace VideoStream {
  export type AsObject = {
    index: number,
    codecName: string,
    codecLongName: string,
    profile: string,
    codecType: string,
    width: number,
    height: number,
    rFrameRate: string,
  }
}

export class VideoFormat extends jspb.Message {
  getFilename(): string;
  setFilename(value: string): VideoFormat;

  getNbStreams(): number;
  setNbStreams(value: number): VideoFormat;

  getFormatName(): string;
  setFormatName(value: string): VideoFormat;

  getFormatLongName(): string;
  setFormatLongName(value: string): VideoFormat;

  getStartTime(): string;
  setStartTime(value: string): VideoFormat;

  getDuration(): string;
  setDuration(value: string): VideoFormat;

  getSize(): string;
  setSize(value: string): VideoFormat;

  getBitRate(): string;
  setBitRate(value: string): VideoFormat;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VideoFormat.AsObject;
  static toObject(includeInstance: boolean, msg: VideoFormat): VideoFormat.AsObject;
  static serializeBinaryToWriter(message: VideoFormat, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VideoFormat;
  static deserializeBinaryFromReader(message: VideoFormat, reader: jspb.BinaryReader): VideoFormat;
}

export namespace VideoFormat {
  export type AsObject = {
    filename: string,
    nbStreams: number,
    formatName: string,
    formatLongName: string,
    startTime: string,
    duration: string,
    size: string,
    bitRate: string,
  }
}

export class VideoMetadata extends jspb.Message {
  getStreamsList(): Array<VideoStream>;
  setStreamsList(value: Array<VideoStream>): VideoMetadata;
  clearStreamsList(): VideoMetadata;
  addStreams(value?: VideoStream, index?: number): VideoStream;

  getFormat(): VideoFormat | undefined;
  setFormat(value?: VideoFormat): VideoMetadata;
  hasFormat(): boolean;
  clearFormat(): VideoMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VideoMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: VideoMetadata): VideoMetadata.AsObject;
  static serializeBinaryToWriter(message: VideoMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VideoMetadata;
  static deserializeBinaryFromReader(message: VideoMetadata, reader: jspb.BinaryReader): VideoMetadata;
}

export namespace VideoMetadata {
  export type AsObject = {
    streamsList: Array<VideoStream.AsObject>,
    format?: VideoFormat.AsObject,
  }
}

export class Video extends jspb.Message {
  getId(): number;
  setId(value: number): Video;

  getName(): string;
  setName(value: string): Video;

  getMeta(): VideoMetadata | undefined;
  setMeta(value?: VideoMetadata): Video;
  hasMeta(): boolean;
  clearMeta(): Video;

  getSubtitlePathsList(): Array<string>;
  setSubtitlePathsList(value: Array<string>): Video;
  clearSubtitlePathsList(): Video;
  addSubtitlePaths(value: string, index?: number): Video;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Video.AsObject;
  static toObject(includeInstance: boolean, msg: Video): Video.AsObject;
  static serializeBinaryToWriter(message: Video, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Video;
  static deserializeBinaryFromReader(message: Video, reader: jspb.BinaryReader): Video;
}

export namespace Video {
  export type AsObject = {
    id: number,
    name: string,
    meta?: VideoMetadata.AsObject,
    subtitlePathsList: Array<string>,
  }
}

