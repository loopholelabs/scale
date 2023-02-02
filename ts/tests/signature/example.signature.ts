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

import {
    encodeString,
    decodeString,
    encodeError,
} from "@loopholelabs/polyglot-ts";

export class ExampleContext {
    constructor(data: string) {
        this._Data = data
    }

    private _Data: string;

    get Data(): string {
        return this._Data
    }

    set Data(data: string) {
        this._Data = data
    }

    encode(buf: Uint8Array): Uint8Array {
        let encoded = buf
        encoded = encodeString(encoded, this._Data)
        return encoded
    }

    internalError(buf: Uint8Array, err: Error): Uint8Array {
        return encodeError(buf, err)
    }

    static decode(buf: Uint8Array): {
        buf: Uint8Array,
        value: ExampleContext
    } {
        let decoded = buf
        const Data = decodeString(decoded)
        decoded = Data.buf
        return {
            buf: decoded,
            value: new ExampleContext(Data.value)
        }
    }
}