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

export class HashSlotAction extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): HashSlotAction;

  getStep(): number;
  setStep(value: number): HashSlotAction;

  getSlots(): Uint8Array | string;
  getSlots_asU8(): Uint8Array;
  getSlots_asB64(): string;
  setSlots(value: Uint8Array | string): HashSlotAction;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HashSlotAction.AsObject;
  static toObject(includeInstance: boolean, msg: HashSlotAction): HashSlotAction.AsObject;
  static serializeBinaryToWriter(message: HashSlotAction, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HashSlotAction;
  static deserializeBinaryFromReader(message: HashSlotAction, reader: jspb.BinaryReader): HashSlotAction;
}

export namespace HashSlotAction {
  export type AsObject = {
    myId: string,
    step: number,
    slots: Uint8Array | string,
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

  getHashSlot(): HashSlotAction | undefined;
  setHashSlot(value?: HashSlotAction): RaftTransaction;
  hasHashSlot(): boolean;
  clearHashSlot(): RaftTransaction;

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
    hashSlot?: HashSlotAction.AsObject,
  }

  export enum Type { 
    UNKNOWN = 0,
    NEWNODE = 1,
    HASHSLOTACTION = 2,
  }
}

export class NodeState extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): NodeState;

  getState(): number;
  setState(value: number): NodeState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NodeState.AsObject;
  static toObject(includeInstance: boolean, msg: NodeState): NodeState.AsObject;
  static serializeBinaryToWriter(message: NodeState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NodeState;
  static deserializeBinaryFromReader(message: NodeState, reader: jspb.BinaryReader): NodeState;
}

export namespace NodeState {
  export type AsObject = {
    myId: string,
    state: number,
  }
}

export class RaftPing extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftPing;

  getRole(): number;
  setRole(value: number): RaftPing;

  getCurrentEpoch(): number;
  setCurrentEpoch(value: number): RaftPing;

  getCommitedEpoch(): number;
  setCommitedEpoch(value: number): RaftPing;

  getVoteEpoch(): number;
  setVoteEpoch(value: number): RaftPing;

  getNodeStatesList(): Array<NodeState>;
  setNodeStatesList(value: Array<NodeState>): RaftPing;
  clearNodeStatesList(): RaftPing;
  addNodeStates(value?: NodeState, index?: number): NodeState;

  getMasterAddress(): string;
  setMasterAddress(value: string): RaftPing;

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
    currentEpoch: number,
    commitedEpoch: number,
    voteEpoch: number,
    nodeStatesList: Array<NodeState.AsObject>,
    masterAddress: string,
  }
}

export class RaftPong extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftPong;

  getRole(): number;
  setRole(value: number): RaftPong;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftPong.AsObject;
  static toObject(includeInstance: boolean, msg: RaftPong): RaftPong.AsObject;
  static serializeBinaryToWriter(message: RaftPong, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftPong;
  static deserializeBinaryFromReader(message: RaftPong, reader: jspb.BinaryReader): RaftPong;
}

export namespace RaftPong {
  export type AsObject = {
    myId: string,
    role: number,
  }
}

export class RaftElection extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftElection;

  getVoteEpoch(): number;
  setVoteEpoch(value: number): RaftElection;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftElection.AsObject;
  static toObject(includeInstance: boolean, msg: RaftElection): RaftElection.AsObject;
  static serializeBinaryToWriter(message: RaftElection, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftElection;
  static deserializeBinaryFromReader(message: RaftElection, reader: jspb.BinaryReader): RaftElection;
}

export namespace RaftElection {
  export type AsObject = {
    myId: string,
    voteEpoch: number,
  }
}

export class RaftElectionRet extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftElectionRet;

  getGotVoteId(): string;
  setGotVoteId(value: string): RaftElectionRet;

  getVoteEpoch(): number;
  setVoteEpoch(value: number): RaftElectionRet;

  getSuccess(): boolean;
  setSuccess(value: boolean): RaftElectionRet;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftElectionRet.AsObject;
  static toObject(includeInstance: boolean, msg: RaftElectionRet): RaftElectionRet.AsObject;
  static serializeBinaryToWriter(message: RaftElectionRet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftElectionRet;
  static deserializeBinaryFromReader(message: RaftElectionRet, reader: jspb.BinaryReader): RaftElectionRet;
}

export namespace RaftElectionRet {
  export type AsObject = {
    myId: string,
    gotVoteId: string,
    voteEpoch: number,
    success: boolean,
  }
}

export class RaftSyncActions extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftSyncActions;

  getCurrentEpoch(): number;
  setCurrentEpoch(value: number): RaftSyncActions;

  getCommitedEpoch(): number;
  setCommitedEpoch(value: number): RaftSyncActions;

  getActionsList(): Array<RaftTransaction>;
  setActionsList(value: Array<RaftTransaction>): RaftSyncActions;
  clearActionsList(): RaftSyncActions;
  addActions(value?: RaftTransaction, index?: number): RaftTransaction;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftSyncActions.AsObject;
  static toObject(includeInstance: boolean, msg: RaftSyncActions): RaftSyncActions.AsObject;
  static serializeBinaryToWriter(message: RaftSyncActions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftSyncActions;
  static deserializeBinaryFromReader(message: RaftSyncActions, reader: jspb.BinaryReader): RaftSyncActions;
}

export namespace RaftSyncActions {
  export type AsObject = {
    myId: string,
    currentEpoch: number,
    commitedEpoch: number,
    actionsList: Array<RaftTransaction.AsObject>,
  }
}

export class RaftSyncActionsRet extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftSyncActionsRet;

  getCurrentEpoch(): number;
  setCurrentEpoch(value: number): RaftSyncActionsRet;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftSyncActionsRet.AsObject;
  static toObject(includeInstance: boolean, msg: RaftSyncActionsRet): RaftSyncActionsRet.AsObject;
  static serializeBinaryToWriter(message: RaftSyncActionsRet, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftSyncActionsRet;
  static deserializeBinaryFromReader(message: RaftSyncActionsRet, reader: jspb.BinaryReader): RaftSyncActionsRet;
}

export namespace RaftSyncActionsRet {
  export type AsObject = {
    myId: string,
    currentEpoch: number,
  }
}

export class RaftReqActions extends jspb.Message {
  getMyId(): string;
  setMyId(value: string): RaftReqActions;

  getCommitedEpoch(): number;
  setCommitedEpoch(value: number): RaftReqActions;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RaftReqActions.AsObject;
  static toObject(includeInstance: boolean, msg: RaftReqActions): RaftReqActions.AsObject;
  static serializeBinaryToWriter(message: RaftReqActions, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RaftReqActions;
  static deserializeBinaryFromReader(message: RaftReqActions, reader: jspb.BinaryReader): RaftReqActions;
}

export namespace RaftReqActions {
  export type AsObject = {
    myId: string,
    commitedEpoch: number,
  }
}

export class SlotMsg extends jspb.Message {
  getSlot(): number;
  setSlot(value: number): SlotMsg;

  getMsg(): string;
  setMsg(value: string): SlotMsg;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SlotMsg.AsObject;
  static toObject(includeInstance: boolean, msg: SlotMsg): SlotMsg.AsObject;
  static serializeBinaryToWriter(message: SlotMsg, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SlotMsg;
  static deserializeBinaryFromReader(message: SlotMsg, reader: jspb.BinaryReader): SlotMsg;
}

export namespace SlotMsg {
  export type AsObject = {
    slot: number,
    msg: string,
  }
}

export class RaftMsg extends jspb.Message {
  getType(): RaftMsg.Type;
  setType(value: RaftMsg.Type): RaftMsg;

  getAction(): RaftTransaction | undefined;
  setAction(value?: RaftTransaction): RaftMsg;
  hasAction(): boolean;
  clearAction(): RaftMsg;

  getSyncActions(): RaftSyncActions | undefined;
  setSyncActions(value?: RaftSyncActions): RaftMsg;
  hasSyncActions(): boolean;
  clearSyncActions(): RaftMsg;

  getSyncActionsRet(): RaftSyncActionsRet | undefined;
  setSyncActionsRet(value?: RaftSyncActionsRet): RaftMsg;
  hasSyncActionsRet(): boolean;
  clearSyncActionsRet(): RaftMsg;

  getPing(): RaftPing | undefined;
  setPing(value?: RaftPing): RaftMsg;
  hasPing(): boolean;
  clearPing(): RaftMsg;

  getPong(): RaftPong | undefined;
  setPong(value?: RaftPong): RaftMsg;
  hasPong(): boolean;
  clearPong(): RaftMsg;

  getElection(): RaftElection | undefined;
  setElection(value?: RaftElection): RaftMsg;
  hasElection(): boolean;
  clearElection(): RaftMsg;

  getElectionRet(): RaftElectionRet | undefined;
  setElectionRet(value?: RaftElectionRet): RaftMsg;
  hasElectionRet(): boolean;
  clearElectionRet(): RaftMsg;

  getReqActions(): RaftReqActions | undefined;
  setReqActions(value?: RaftReqActions): RaftMsg;
  hasReqActions(): boolean;
  clearReqActions(): RaftMsg;

  getSlotMsg(): SlotMsg | undefined;
  setSlotMsg(value?: SlotMsg): RaftMsg;
  hasSlotMsg(): boolean;
  clearSlotMsg(): RaftMsg;

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
    syncActions?: RaftSyncActions.AsObject,
    syncActionsRet?: RaftSyncActionsRet.AsObject,
    ping?: RaftPing.AsObject,
    pong?: RaftPong.AsObject,
    election?: RaftElection.AsObject,
    electionRet?: RaftElectionRet.AsObject,
    reqActions?: RaftReqActions.AsObject,
    slotMsg?: SlotMsg.AsObject,
  }

  export enum Type { 
    UNKNOWN = 0,
    ACTION = 1,
    SYNCACTION = 2,
    SYNCACTIONRET = 3,
    PING = 4,
    PONG = 5,
    ELECTION = 6,
    ELECTIONRET = 7,
    REQACTIONS = 8,
    SENDMSG2SLOT = 9,
  }
}

