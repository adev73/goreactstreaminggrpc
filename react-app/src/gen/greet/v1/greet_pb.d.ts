// @generated by protoc-gen-es v1.2.1 with parameter "target=js+dts"
// @generated from file greet/v1/greet.proto (package greet.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message greet.v1.GreetRequest
 */
export declare class GreetRequest extends Message<GreetRequest> {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId: string;

  /**
   * @generated from field: string name = 2;
   */
  name: string;

  /**
   * @generated from field: bool end_session = 3;
   */
  endSession: boolean;

  constructor(data?: PartialMessage<GreetRequest>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "greet.v1.GreetRequest";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GreetRequest;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GreetRequest;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GreetRequest;

  static equals(a: GreetRequest | PlainMessage<GreetRequest> | undefined, b: GreetRequest | PlainMessage<GreetRequest> | undefined): boolean;
}

/**
 * @generated from message greet.v1.GreetResponse
 */
export declare class GreetResponse extends Message<GreetResponse> {
  /**
   * @generated from field: bool confirmed = 2;
   */
  confirmed: boolean;

  constructor(data?: PartialMessage<GreetResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "greet.v1.GreetResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GreetResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GreetResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GreetResponse;

  static equals(a: GreetResponse | PlainMessage<GreetResponse> | undefined, b: GreetResponse | PlainMessage<GreetResponse> | undefined): boolean;
}

/**
 * @generated from message greet.v1.GreetingsRequest
 */
export declare class GreetingsRequest extends Message<GreetingsRequest> {
  constructor(data?: PartialMessage<GreetingsRequest>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "greet.v1.GreetingsRequest";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GreetingsRequest;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GreetingsRequest;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GreetingsRequest;

  static equals(a: GreetingsRequest | PlainMessage<GreetingsRequest> | undefined, b: GreetingsRequest | PlainMessage<GreetingsRequest> | undefined): boolean;
}

/**
 * @generated from message greet.v1.GreetingsResponse
 */
export declare class GreetingsResponse extends Message<GreetingsResponse> {
  /**
   * @generated from field: string session_id = 1;
   */
  sessionId: string;

  /**
   * @generated from field: string greeting = 2;
   */
  greeting: string;

  /**
   * @generated from field: bool end_session = 3;
   */
  endSession: boolean;

  constructor(data?: PartialMessage<GreetingsResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "greet.v1.GreetingsResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GreetingsResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GreetingsResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GreetingsResponse;

  static equals(a: GreetingsResponse | PlainMessage<GreetingsResponse> | undefined, b: GreetingsResponse | PlainMessage<GreetingsResponse> | undefined): boolean;
}
