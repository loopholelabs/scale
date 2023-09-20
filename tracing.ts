/*
        Copyright 2022 Loophole Labs

        Licensed under the Apache License, Version 2.0 (the "License");
        you may not use this file except in compliance with the License.
        You may obtain a copy of the License at

                   http://www.apache.org/licenses/LICENSE-2.0

        Unless required by applicable law or agreed to in writing, software
        distributed under the License is distributed on an "AS IS" BASIS,
        WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
        See the License for the specific language governing permissions and
        limitations under the License.
*/

import {Module} from "./module";
import {Signature} from "@loopholelabs/scale-signature-interfaces";

const ErrNoInstance = Error("no webassembly instance")
const ErrNoMemory = Error("no exported memory in webassembly instance")

export interface HostFunctions extends WebAssembly.ModuleImports {
  get_function_name_len(): number
  get_function_name(ptr: number): void
  get_instance_id(ptr: number): void
  otel_trace_json(ptr: number, len: number): void
}

export type CallbackFunction = (data: string) => void

export class Tracing<T extends Signature> {
  private exports: WebAssembly.Exports | undefined;

  private readonly module: Module<T>;
  private readonly traceCallback: CallbackFunction | undefined;

  constructor(module: Module<T>, traceCallback: CallbackFunction | undefined) {
    this.module = module;
    this.traceCallback = traceCallback;
  }

  getFunctionNameLen(): number {
    const enc = new TextEncoder();
    const data = enc.encode(this.module.template.identifier);
    return data.length;
  }

  getFunctionName(ptr: number) {
    const enc = new TextEncoder();
    const data = enc.encode(this.module.template.identifier);
    const buffer = this.getDataView();
    for (let i=0;i<data.length;i++) {
      const d = data.at(i);
      if (typeof d !== "undefined") {
        buffer.setInt8(ptr+i, d);
      }
    }
  }

  getInstanceId(ptr: number) {
    if (typeof this.module.function === "undefined") return;
    const buffer = this.getDataView();
    for (let i=0;i<this.module.function.instance.identifier.length;i++) {
      const d = this.module.function.instance.identifier.at(i);
      if (typeof d !== "undefined") {
        buffer.setInt8(ptr+i, d);
      }
    }
  }

  otelTraceJSON(ptr: number, len: number) {
    if (typeof this.traceCallback === "undefined") return;
    const buffer = this.getDataView();
    const data = buffer.buffer.slice(ptr, ptr + len);
    const dec = new TextDecoder();
    const s = dec.decode(data);
    this.traceCallback(s);
  }

  GetImports(): HostFunctions {
    return {
      get_instance_id: this.getInstanceId.bind(this),
      get_function_name_len: this.getFunctionNameLen.bind(this),
      get_function_name: this.getFunctionName.bind(this),
      otel_trace_json: this.otelTraceJSON.bind(this),
    }
  }

  private getDataView(): DataView {
    if (!this.exports) {
        throw ErrNoInstance;
    }
    if (!this.exports.memory) {
        throw ErrNoMemory;
    }

    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    return new DataView(this.exports.memory.buffer);
  }
}