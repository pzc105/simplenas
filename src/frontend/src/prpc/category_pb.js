// source: category.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global = (function() {
  if (this) { return this; }
  if (typeof window !== 'undefined') { return window; }
  if (typeof global !== 'undefined') { return global; }
  if (typeof self !== 'undefined') { return self; }
  return Function('return this')();
}.call(null));

goog.exportSymbol('proto.prpc.CategoryItem', null, global);
goog.exportSymbol('proto.prpc.CategoryItem.Type', null, global);
goog.exportSymbol('proto.prpc.SharedItem', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.prpc.CategoryItem = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, proto.prpc.CategoryItem.repeatedFields_, null);
};
goog.inherits(proto.prpc.CategoryItem, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.prpc.CategoryItem.displayName = 'proto.prpc.CategoryItem';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.prpc.SharedItem = function(opt_data) {
  jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.prpc.SharedItem, jspb.Message);
if (goog.DEBUG && !COMPILED) {
  /**
   * @public
   * @override
   */
  proto.prpc.SharedItem.displayName = 'proto.prpc.SharedItem';
}

/**
 * List of repeated fields within this message type.
 * @private {!Array<number>}
 * @const
 */
proto.prpc.CategoryItem.repeatedFields_ = [10];



if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.prpc.CategoryItem.prototype.toObject = function(opt_includeInstance) {
  return proto.prpc.CategoryItem.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.prpc.CategoryItem} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.prpc.CategoryItem.toObject = function(includeInstance, msg) {
  var f, obj = {
    id: jspb.Message.getFieldWithDefault(msg, 1, 0),
    typeId: jspb.Message.getFieldWithDefault(msg, 2, 0),
    creator: jspb.Message.getFieldWithDefault(msg, 3, 0),
    name: jspb.Message.getFieldWithDefault(msg, 4, ""),
    resourcePath: jspb.Message.getFieldWithDefault(msg, 5, ""),
    posterPath: jspb.Message.getFieldWithDefault(msg, 6, ""),
    introduce: jspb.Message.getFieldWithDefault(msg, 7, ""),
    other: jspb.Message.getFieldWithDefault(msg, 8, ""),
    parentId: jspb.Message.getFieldWithDefault(msg, 9, 0),
    subItemIdsList: (f = jspb.Message.getRepeatedField(msg, 10)) == null ? undefined : f
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.prpc.CategoryItem}
 */
proto.prpc.CategoryItem.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.prpc.CategoryItem;
  return proto.prpc.CategoryItem.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.prpc.CategoryItem} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.prpc.CategoryItem}
 */
proto.prpc.CategoryItem.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setId(value);
      break;
    case 2:
      var value = /** @type {!proto.prpc.CategoryItem.Type} */ (reader.readEnum());
      msg.setTypeId(value);
      break;
    case 3:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setCreator(value);
      break;
    case 4:
      var value = /** @type {string} */ (reader.readString());
      msg.setName(value);
      break;
    case 5:
      var value = /** @type {string} */ (reader.readString());
      msg.setResourcePath(value);
      break;
    case 6:
      var value = /** @type {string} */ (reader.readString());
      msg.setPosterPath(value);
      break;
    case 7:
      var value = /** @type {string} */ (reader.readString());
      msg.setIntroduce(value);
      break;
    case 8:
      var value = /** @type {string} */ (reader.readString());
      msg.setOther(value);
      break;
    case 9:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setParentId(value);
      break;
    case 10:
      var values = /** @type {!Array<number>} */ (reader.isDelimited() ? reader.readPackedInt64() : [reader.readInt64()]);
      for (var i = 0; i < values.length; i++) {
        msg.addSubItemIds(values[i]);
      }
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.prpc.CategoryItem.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.prpc.CategoryItem.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.prpc.CategoryItem} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.prpc.CategoryItem.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getId();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getTypeId();
  if (f !== 0.0) {
    writer.writeEnum(
      2,
      f
    );
  }
  f = message.getCreator();
  if (f !== 0) {
    writer.writeInt64(
      3,
      f
    );
  }
  f = message.getName();
  if (f.length > 0) {
    writer.writeString(
      4,
      f
    );
  }
  f = message.getResourcePath();
  if (f.length > 0) {
    writer.writeString(
      5,
      f
    );
  }
  f = message.getPosterPath();
  if (f.length > 0) {
    writer.writeString(
      6,
      f
    );
  }
  f = message.getIntroduce();
  if (f.length > 0) {
    writer.writeString(
      7,
      f
    );
  }
  f = message.getOther();
  if (f.length > 0) {
    writer.writeString(
      8,
      f
    );
  }
  f = message.getParentId();
  if (f !== 0) {
    writer.writeInt64(
      9,
      f
    );
  }
  f = message.getSubItemIdsList();
  if (f.length > 0) {
    writer.writePackedInt64(
      10,
      f
    );
  }
};


/**
 * @enum {number}
 */
proto.prpc.CategoryItem.Type = {
  UNKNOWN: 0,
  HOME: 1,
  DIRECTORY: 2,
  VIDEO: 3,
  OTHER: 4,
  AUDIO: 5,
  MAGNETURI: 6
};

/**
 * optional int64 id = 1;
 * @return {number}
 */
proto.prpc.CategoryItem.prototype.getId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional Type type_id = 2;
 * @return {!proto.prpc.CategoryItem.Type}
 */
proto.prpc.CategoryItem.prototype.getTypeId = function() {
  return /** @type {!proto.prpc.CategoryItem.Type} */ (jspb.Message.getFieldWithDefault(this, 2, 0));
};


/**
 * @param {!proto.prpc.CategoryItem.Type} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setTypeId = function(value) {
  return jspb.Message.setProto3EnumField(this, 2, value);
};


/**
 * optional int64 creator = 3;
 * @return {number}
 */
proto.prpc.CategoryItem.prototype.getCreator = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 3, 0));
};


/**
 * @param {number} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setCreator = function(value) {
  return jspb.Message.setProto3IntField(this, 3, value);
};


/**
 * optional string name = 4;
 * @return {string}
 */
proto.prpc.CategoryItem.prototype.getName = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setName = function(value) {
  return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional string resource_path = 5;
 * @return {string}
 */
proto.prpc.CategoryItem.prototype.getResourcePath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 5, ""));
};


/**
 * @param {string} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setResourcePath = function(value) {
  return jspb.Message.setProto3StringField(this, 5, value);
};


/**
 * optional string poster_path = 6;
 * @return {string}
 */
proto.prpc.CategoryItem.prototype.getPosterPath = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 6, ""));
};


/**
 * @param {string} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setPosterPath = function(value) {
  return jspb.Message.setProto3StringField(this, 6, value);
};


/**
 * optional string introduce = 7;
 * @return {string}
 */
proto.prpc.CategoryItem.prototype.getIntroduce = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 7, ""));
};


/**
 * @param {string} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setIntroduce = function(value) {
  return jspb.Message.setProto3StringField(this, 7, value);
};


/**
 * optional string other = 8;
 * @return {string}
 */
proto.prpc.CategoryItem.prototype.getOther = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 8, ""));
};


/**
 * @param {string} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setOther = function(value) {
  return jspb.Message.setProto3StringField(this, 8, value);
};


/**
 * optional int64 parent_id = 9;
 * @return {number}
 */
proto.prpc.CategoryItem.prototype.getParentId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 9, 0));
};


/**
 * @param {number} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setParentId = function(value) {
  return jspb.Message.setProto3IntField(this, 9, value);
};


/**
 * repeated int64 sub_item_ids = 10;
 * @return {!Array<number>}
 */
proto.prpc.CategoryItem.prototype.getSubItemIdsList = function() {
  return /** @type {!Array<number>} */ (jspb.Message.getRepeatedField(this, 10));
};


/**
 * @param {!Array<number>} value
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.setSubItemIdsList = function(value) {
  return jspb.Message.setField(this, 10, value || []);
};


/**
 * @param {number} value
 * @param {number=} opt_index
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.addSubItemIds = function(value, opt_index) {
  return jspb.Message.addToRepeatedField(this, 10, value, opt_index);
};


/**
 * Clears the list making it empty but non-null.
 * @return {!proto.prpc.CategoryItem} returns this
 */
proto.prpc.CategoryItem.prototype.clearSubItemIdsList = function() {
  return this.setSubItemIdsList([]);
};





if (jspb.Message.GENERATE_TO_OBJECT) {
/**
 * Creates an object representation of this proto.
 * Field names that are reserved in JavaScript and will be renamed to pb_name.
 * Optional fields that are not set will be set to undefined.
 * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
 * For the list of reserved names please see:
 *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
 * @param {boolean=} opt_includeInstance Deprecated. whether to include the
 *     JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @return {!Object}
 */
proto.prpc.SharedItem.prototype.toObject = function(opt_includeInstance) {
  return proto.prpc.SharedItem.toObject(opt_includeInstance, this);
};


/**
 * Static version of the {@see toObject} method.
 * @param {boolean|undefined} includeInstance Deprecated. Whether to include
 *     the JSPB instance for transitional soy proto support:
 *     http://goto/soy-param-migration
 * @param {!proto.prpc.SharedItem} msg The msg instance to transform.
 * @return {!Object}
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.prpc.SharedItem.toObject = function(includeInstance, msg) {
  var f, obj = {
    itemId: jspb.Message.getFieldWithDefault(msg, 1, 0),
    shareId: jspb.Message.getFieldWithDefault(msg, 2, "")
  };

  if (includeInstance) {
    obj.$jspbMessageInstance = msg;
  }
  return obj;
};
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.prpc.SharedItem}
 */
proto.prpc.SharedItem.deserializeBinary = function(bytes) {
  var reader = new jspb.BinaryReader(bytes);
  var msg = new proto.prpc.SharedItem;
  return proto.prpc.SharedItem.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.prpc.SharedItem} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.prpc.SharedItem}
 */
proto.prpc.SharedItem.deserializeBinaryFromReader = function(msg, reader) {
  while (reader.nextField()) {
    if (reader.isEndGroup()) {
      break;
    }
    var field = reader.getFieldNumber();
    switch (field) {
    case 1:
      var value = /** @type {number} */ (reader.readInt64());
      msg.setItemId(value);
      break;
    case 2:
      var value = /** @type {string} */ (reader.readString());
      msg.setShareId(value);
      break;
    default:
      reader.skipField();
      break;
    }
  }
  return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.prpc.SharedItem.prototype.serializeBinary = function() {
  var writer = new jspb.BinaryWriter();
  proto.prpc.SharedItem.serializeBinaryToWriter(this, writer);
  return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.prpc.SharedItem} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.prpc.SharedItem.serializeBinaryToWriter = function(message, writer) {
  var f = undefined;
  f = message.getItemId();
  if (f !== 0) {
    writer.writeInt64(
      1,
      f
    );
  }
  f = message.getShareId();
  if (f.length > 0) {
    writer.writeString(
      2,
      f
    );
  }
};


/**
 * optional int64 item_id = 1;
 * @return {number}
 */
proto.prpc.SharedItem.prototype.getItemId = function() {
  return /** @type {number} */ (jspb.Message.getFieldWithDefault(this, 1, 0));
};


/**
 * @param {number} value
 * @return {!proto.prpc.SharedItem} returns this
 */
proto.prpc.SharedItem.prototype.setItemId = function(value) {
  return jspb.Message.setProto3IntField(this, 1, value);
};


/**
 * optional string share_id = 2;
 * @return {string}
 */
proto.prpc.SharedItem.prototype.getShareId = function() {
  return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.prpc.SharedItem} returns this
 */
proto.prpc.SharedItem.prototype.setShareId = function(value) {
  return jspb.Message.setProto3StringField(this, 2, value);
};


goog.object.extend(exports, proto.prpc);
