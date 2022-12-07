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

import { Signature } from "../signature/signature";

import { SFunction } from "./sfunction";
import { Runtime } from "./runtime";
import { Module } from "./module";

export class Pool<T extends Signature> {
  private f: SFunction<T>;
  private r: Runtime<T>;

  constructor(f: SFunction<T>, r: Runtime<T>) {
    this.f = f;
    this.r = r;
  }

  // For now, we don't actually have a pool, we just create new Modules each time.
  Get(): Module<T> {
    return new Module<T>(this.f, this.r);
  }
}