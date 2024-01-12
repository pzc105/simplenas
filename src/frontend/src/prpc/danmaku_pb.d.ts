import * as jspb from 'google-protobuf'



export class Danmaku extends jspb.Message {
  getId(): string;
  setId(value: string): Danmaku;

  getUserId(): number;
  setUserId(value: number): Danmaku;

  getUserName(): string;
  setUserName(value: string): Danmaku;

  getSTime(): number;
  setSTime(value: number): Danmaku;

  getText(): string;
  setText(value: string): Danmaku;

  getType(): number;
  setType(value: number): Danmaku;

  getColor(): number;
  setColor(value: number): Danmaku;

  getDTime(): number;
  setDTime(value: number): Danmaku;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Danmaku.AsObject;
  static toObject(includeInstance: boolean, msg: Danmaku): Danmaku.AsObject;
  static serializeBinaryToWriter(message: Danmaku, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Danmaku;
  static deserializeBinaryFromReader(message: Danmaku, reader: jspb.BinaryReader): Danmaku;
}

export namespace Danmaku {
  export type AsObject = {
    id: string,
    userId: number,
    userName: string,
    sTime: number,
    text: string,
    type: number,
    color: number,
    dTime: number,
  }
}

