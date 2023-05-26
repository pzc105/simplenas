import * as jspb from 'google-protobuf'



export class CategoryItem extends jspb.Message {
  getId(): number;
  setId(value: number): CategoryItem;

  getTypeId(): CategoryItem.Type;
  setTypeId(value: CategoryItem.Type): CategoryItem;

  getCreator(): number;
  setCreator(value: number): CategoryItem;

  getName(): string;
  setName(value: string): CategoryItem;

  getResourcePath(): string;
  setResourcePath(value: string): CategoryItem;

  getPosterPath(): string;
  setPosterPath(value: string): CategoryItem;

  getIntroduce(): string;
  setIntroduce(value: string): CategoryItem;

  getParentId(): number;
  setParentId(value: number): CategoryItem;

  getSubItemIdsList(): Array<number>;
  setSubItemIdsList(value: Array<number>): CategoryItem;
  clearSubItemIdsList(): CategoryItem;
  addSubItemIds(value: number, index?: number): CategoryItem;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CategoryItem.AsObject;
  static toObject(includeInstance: boolean, msg: CategoryItem): CategoryItem.AsObject;
  static serializeBinaryToWriter(message: CategoryItem, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CategoryItem;
  static deserializeBinaryFromReader(message: CategoryItem, reader: jspb.BinaryReader): CategoryItem;
}

export namespace CategoryItem {
  export type AsObject = {
    id: number,
    typeId: CategoryItem.Type,
    creator: number,
    name: string,
    resourcePath: string,
    posterPath: string,
    introduce: string,
    parentId: number,
    subItemIdsList: Array<number>,
  }

  export enum Type { 
    UNKNOWN = 0,
    HOME = 1,
    DIRECTORY = 2,
    VIDEO = 3,
    OTHERFILE = 4,
  }
}

export class SharedItem extends jspb.Message {
  getItemId(): number;
  setItemId(value: number): SharedItem;

  getShareId(): string;
  setShareId(value: string): SharedItem;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SharedItem.AsObject;
  static toObject(includeInstance: boolean, msg: SharedItem): SharedItem.AsObject;
  static serializeBinaryToWriter(message: SharedItem, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SharedItem;
  static deserializeBinaryFromReader(message: SharedItem, reader: jspb.BinaryReader): SharedItem;
}

export namespace SharedItem {
  export type AsObject = {
    itemId: number,
    shareId: string,
  }
}

