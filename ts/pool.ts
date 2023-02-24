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

import { Signature } from "@loopholelabs/scale-signature";

import { Func } from "./func";
import { RawRuntime } from "./runtime";
import { Module } from "./module";

export class Pool<T extends Signature> {
  private readonly f: Func<T>;
  private readonly r: RawRuntime<T>;

  constructor(f: Func<T>, r: RawRuntime<T>) {
    this.f = f;
    this.r = r;
  }

  Get(): Module<T> {
    return new Module<T>(this.f, this.r);
  }
}