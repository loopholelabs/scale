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

const ErrNoInstance = Error("no webassembly instance")
const ErrNoMemory = Error("no exported memory in webassembly instance")

export interface ScaleHostFuncs extends WebAssembly.ModuleImports {
  get_invocation_id(ptr: number): void
  get_service_name_len(): number
  get_service_name(ptr: number): void
  send_otel_trace_json(ptr: number, len: number): void
}

export type TraceCallbackFunc = (data: string) => void

export class ScaleMod {
  private exports: WebAssembly.Exports | undefined;

  private invocationId: Buffer;
  private serviceName: string;
  private traceCallback: TraceCallbackFunc | undefined;

  constructor(serviceName: string, invocationId: Buffer, traceCallback: TraceCallbackFunc | undefined) {
    this.invocationId = invocationId;
    this.serviceName = serviceName;
    this.traceCallback = traceCallback;
  }

  getInvocationId(ptr: number) {
    let buffer = this.getDataView();
    for (let i=0;i<this.invocationId.length;i++) {
      let d = this.invocationId.at(i);
      if (d!=undefined) {
        buffer.setInt8(ptr+i, d);
      }
    }
  }

  getServiceNameLen(): number {
    let enc = new TextEncoder();
    let data = enc.encode(this.serviceName);
    return data.length;
  }

  getServiceName(ptr: number) {
    let enc = new TextEncoder();
    let data = enc.encode(this.serviceName);
    let buffer = this.getDataView();
    for (let i=0;i<data.length;i++) {
      let d = data.at(i);
      if (d!=undefined) {
        buffer.setInt8(ptr+i, d);
      }
    }
  }

  sendOtelTraceJson(ptr: number, len: number) {
    if (this.traceCallback==undefined) return;
    let buffer = this.getDataView();
    let data = buffer.buffer.slice(ptr, ptr + len);
    let dec = new TextDecoder();
    let s = dec.decode(data);
    this.traceCallback(s);
  }

  GetImports(): ScaleHostFuncs {
    return {
      get_invocation_id: this.getInvocationId.bind(this),
      get_service_name_len: this.getServiceNameLen.bind(this),
      get_service_name: this.getServiceName.bind(this),
      send_otel_trace_json: this.sendOtelTraceJson.bind(this),
    }
  }

  private getDataView(): DataView {
    if (!this.exports) {
        throw ErrNoInstance;
    }
    if (!this.exports.memory) {
        throw ErrNoMemory;
    }
    // @ts-ignore
    return new DataView(this.exports.memory.buffer);
  }

  public SetInstance(instance: WebAssembly.Instance) {
      this.exports = instance.exports;
  }

}