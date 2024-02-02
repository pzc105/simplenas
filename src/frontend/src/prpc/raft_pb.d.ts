import * as jspb from 'google-protobuf'



export class JoinMsg extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): JoinMsg;

  getMyAddress(): string;
  setMyAddress(value: string): JoinMsg;

  getRole(): number;
  setRole(value: number): JoinMsg;

  getSlots(): Uint8Array | string;
  getSlots_asU8(): Uint8Array;
  getSlots_asB64(): string;
  setSlots(value: Uint8Array | string): JoinMsg;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinMsg.AsObject;
  static toObject(includeInstance: boolean, msg: JoinMsg): JoinMsg.AsObject;
  static serializeBinaryToWriter(message: JoinMsg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinMsg;
  static deserializeBinaryFromReader(message: JoinMsg, reader: jspb.BinaryReader): JoinMsg;
}

export namespace JoinMsg {
  export type AsObject = {
    myId: string,
    myAddress: string,
    role: number,
    slots: Uint8Array | string,
  }
}

export class JoinRet extends jspb.Message {
  getSuccess(): boolean;
  setSuccess(value: boolean): JoinRet;

  getCurrentEpoch(): number;
  setCurrentEpoch(value: number): JoinRet;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): JoinRet.AsObject;
  static toObject(includeInstance: boolean, msg: JoinRet): JoinRet.AsObject;
  static serializeBinaryToWriter(message: JoinRet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): JoinRet;
  static deserializeBinaryFromReader(message: JoinRet, reader: jspb.BinaryReader): JoinRet;
}

export namespace JoinRet {
  export type AsObject = {
    success: boolean,
    currentEpoch: number,
  }
}

export class SendMsg extends jspb.Message {
  getType(): SendMsg.Type;
  setType(value: SendMsg.Type): SendMsg;

  getJoinmsg(): JoinMsg | undefined;
  setJoinmsg(value?: JoinMsg): SendMsg;
  hasJoinmsg(): boolean;
  clearJoinmsg(): SendMsg;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendMsg.AsObject;
  static toObject(includeInstance: boolean, msg: SendMsg): SendMsg.AsObject;
  static serializeBinaryToWriter(message: SendMsg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendMsg;
  static deserializeBinaryFromReader(message: SendMsg, reader: jspb.BinaryReader): SendMsg;
}

export namespace SendMsg {
  export type AsObject = {
    type: SendMsg.Type,
    joinmsg?: JoinMsg.AsObject,
  }

  export enum Type { 
    UNKNOWN = 0,
    HELLO = 1,
    JOIN = 2,
  }
}

export class SendRet extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SendRet.AsObject;
  static toObject(includeInstance: boolean, msg: SendRet): SendRet.AsObject;
  static serializeBinaryToWriter(message: SendRet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SendRet;
  static deserializeBinaryFromReader(message: SendRet, reader: jspb.BinaryReader): SendRet;
}

export namespace SendRet {
  export type AsObject = {
  }
}

