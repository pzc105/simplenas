import * as jspb from 'google-protobuf'



export class NewNodeAction extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): NewNodeAction;

  getMyAddress(): string;
  setMyAddress(value: string): NewNodeAction;

  getRole(): number;
  setRole(value: number): NewNodeAction;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NewNodeAction.AsObject;
  static toObject(includeInstance: boolean, msg: NewNodeAction): NewNodeAction.AsObject;
  static serializeBinaryToWriter(message: NewNodeAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NewNodeAction;
  static deserializeBinaryFromReader(message: NewNodeAction, reader: jspb.BinaryReader): NewNodeAction;
}

export namespace NewNodeAction {
  export type AsObject = {
    myId: string,
    myAddress: string,
    role: number,
  }
}

export class RaftTransaction extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftTransaction;

  getType(): RaftTransaction.Type;
  setType(value: RaftTransaction.Type): RaftTransaction;

  getEpoch(): number;
  setEpoch(value: number): RaftTransaction;

  getUserRef(): string;
  setUserRef(value: string): RaftTransaction;

  getNewNode(): NewNodeAction | undefined;
  setNewNode(value?: NewNodeAction): RaftTransaction;
  hasNewNode(): boolean;
  clearNewNode(): RaftTransaction;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftTransaction.AsObject;
  static toObject(includeInstance: boolean, msg: RaftTransaction): RaftTransaction.AsObject;
  static serializeBinaryToWriter(message: RaftTransaction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftTransaction;
  static deserializeBinaryFromReader(message: RaftTransaction, reader: jspb.BinaryReader): RaftTransaction;
}

export namespace RaftTransaction {
  export type AsObject = {
    myId: string,
    type: RaftTransaction.Type,
    epoch: number,
    userRef: string,
    newNode?: NewNodeAction.AsObject,
  }

  export enum Type { 
    UNKNOWN = 0,
    NEWNODE = 1,
  }
}

export class RaftActionRet extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftActionRet;

  getEpoch(): number;
  setEpoch(value: number): RaftActionRet;

  getSuccess(): boolean;
  setSuccess(value: boolean): RaftActionRet;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftActionRet.AsObject;
  static toObject(includeInstance: boolean, msg: RaftActionRet): RaftActionRet.AsObject;
  static serializeBinaryToWriter(message: RaftActionRet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftActionRet;
  static deserializeBinaryFromReader(message: RaftActionRet, reader: jspb.BinaryReader): RaftActionRet;
}

export namespace RaftActionRet {
  export type AsObject = {
    myId: string,
    epoch: number,
    success: boolean,
  }
}

export class RaftPing extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftPing;

  getRole(): number;
  setRole(value: number): RaftPing;

  getEpoch(): number;
  setEpoch(value: number): RaftPing;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftPing.AsObject;
  static toObject(includeInstance: boolean, msg: RaftPing): RaftPing.AsObject;
  static serializeBinaryToWriter(message: RaftPing, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftPing;
  static deserializeBinaryFromReader(message: RaftPing, reader: jspb.BinaryReader): RaftPing;
}

export namespace RaftPing {
  export type AsObject = {
    myId: string,
    role: number,
    epoch: number,
  }
}

export class RaftPong extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftPong.AsObject;
  static toObject(includeInstance: boolean, msg: RaftPong): RaftPong.AsObject;
  static serializeBinaryToWriter(message: RaftPong, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftPong;
  static deserializeBinaryFromReader(message: RaftPong, reader: jspb.BinaryReader): RaftPong;
}

export namespace RaftPong {
  export type AsObject = {
  }
}

export class RaftMsg extends jspb.Message {
  getType(): RaftMsg.Type;
  setType(value: RaftMsg.Type): RaftMsg;

  getAction(): RaftTransaction | undefined;
  setAction(value?: RaftTransaction): RaftMsg;
  hasAction(): boolean;
  clearAction(): RaftMsg;

  getSyncActionsList(): Array<RaftTransaction>;
  setSyncActionsList(value: Array<RaftTransaction>): RaftMsg;
  clearSyncActionsList(): RaftMsg;
  addSyncActions(value?: RaftTransaction, index?: number): RaftTransaction;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftMsg.AsObject;
  static toObject(includeInstance: boolean, msg: RaftMsg): RaftMsg.AsObject;
  static serializeBinaryToWriter(message: RaftMsg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftMsg;
  static deserializeBinaryFromReader(message: RaftMsg, reader: jspb.BinaryReader): RaftMsg;
}

export namespace RaftMsg {
  export type AsObject = {
    type: RaftMsg.Type,
    action?: RaftTransaction.AsObject,
    syncActionsList: Array<RaftTransaction.AsObject>,
  }

  export enum Type { 
    UNKNOWN = 0,
    ACTION = 1,
    SYNCACTION = 2,
    PING = 3,
    PONG = 4,
  }
}

