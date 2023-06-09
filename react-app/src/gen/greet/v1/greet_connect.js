// @generated by protoc-gen-connect-es v0.9.0 with parameter "target=js+dts"
// @generated from file greet/v1/greet.proto (package greet.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { GreetingsRequest, GreetingsResponse, GreetRequest, GreetResponse } from "./greet_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service greet.v1.GreetService
 */
export const GreetService = {
  typeName: "greet.v1.GreetService",
  methods: {
    /**
     * @generated from rpc greet.v1.GreetService.Greet
     */
    greet: {
      name: "Greet",
      I: GreetRequest,
      O: GreetResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc greet.v1.GreetService.Greetings
     */
    greetings: {
      name: "Greetings",
      I: GreetingsRequest,
      O: GreetingsResponse,
      kind: MethodKind.ServerStreaming,
    },
  }
};

