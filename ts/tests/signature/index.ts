import { Kind, decodeError } from "@loopholelabs/polyglot-ts";

import {ExampleContext} from "./example.signature";
import {RuntimeContext as RuntimeContextInterface, Signature} from "@loopholelabs/scale-signature";

export class RuntimeContext implements RuntimeContextInterface {
  private readonly context: ExampleContext;

  constructor(context: ExampleContext) {
    this.context = context;
  }

  Read(b: Uint8Array): Error | undefined {
    if (b.length > 0 && b[0] === Kind.Error) {
      return decodeError(b).value;
    }
    Object.assign(this.context, ExampleContext.decode(b).value);
    return undefined;
  }

  Write(): Uint8Array {
    return this.context.encode(new Uint8Array());
  }

  Error(err: Error): Uint8Array {
    return this.context.internalError(new Uint8Array(), err);
  }
}

export function New(): Context {
  return new Context();
}

export class Context extends ExampleContext implements Signature {
  private readonly runtimeContext;
  constructor() {
    super("");
    this.runtimeContext = new RuntimeContext(this);
  }

  RuntimeContext(): RuntimeContext {
    return this.runtimeContext;
  }
}