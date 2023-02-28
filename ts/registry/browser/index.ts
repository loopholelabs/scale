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

import { ScaleFunc } from "@loopholelabs/scalefile"

let crypto = global.crypto;
if(global.crypto === undefined || global.crypto.subtle === undefined) {
  crypto = require("crypto").webcrypto;
}

import { OpenAPI, RegistryService } from "../../client";

export type PullPolicy = string;
export const AlwaysPullPolicy: PullPolicy = "always";
export const IfNotPresentPullPolicy: PullPolicy = "if-not-present";
export const NeverPullPolicy: PullPolicy = "never";

export const ErrHasMismatch = new Error("hash mismatch");
export const ErrNoFunction = new Error("function does not exist locally and pull policy does not allow pulling from registry");
export const ErrDownloadFailed = new Error("scale function could not be pull from the registry")

export const DefaultBaseURL = "api.scale.sh"
export const DefaultOrganization = "scale"

export interface Config {
  baseURL: string,
  organization: string,
  apiKey: string,
  pullPolicy: PullPolicy,
  cacheDirectory: string,
}

export type Option = (c: Config) => void;

export function WithAPIKey(apiKey: string): Option {
  return (c: Config) => {
    c.apiKey = apiKey;
  }
}

export function WithBaseURL(baseURL: string): Option {
    return (c: Config) => {
        c.baseURL = baseURL;
    }
}

export function WithOrganization(organization: string): Option {
    return (c: Config) => {
        c.organization = organization;
    }
}

export function WithPullPolicy(pullPolicy: PullPolicy): Option {
  return (c: Config) => {
    c.pullPolicy = pullPolicy;
  }
}

export function WithCacheDirectory(cacheDirectory: string): Option {
  return (c: Config) => {
    c.cacheDirectory = cacheDirectory;
  }
}

export async function Download(name: string, tag: string, ...opts: Option[]): Promise<ScaleFunc> {
  const conf = DefaultConfig(...opts);
  const memoryFn = await RegistryService.getRegistryFunction3(conf.organization, name, tag)
  if(!memoryFn.presigned_url || !memoryFn.hash) {
    throw ErrDownloadFailed
  }

  const memoryData = await (await fetch(memoryFn.presigned_url)).arrayBuffer();
  const memoryHash = await ComputeSHA256(memoryData);
  if (memoryHash !== memoryFn.hash) {
    throw ErrHasMismatch;
  }

  return ScaleFunc.Decode(new Uint8Array(memoryData));
}

export function DefaultConfig(...opts: Option[]): Config {
  const conf: Config = {
    baseURL: DefaultBaseURL,
    organization: DefaultOrganization,
    pullPolicy: IfNotPresentPullPolicy,
    apiKey: "",
    cacheDirectory: ""
  }

  for (const opt of opts) {
    opt(conf);
  }

  if (!conf.baseURL || conf.baseURL === "") {
    conf.baseURL = DefaultBaseURL
  }

  if (!conf.organization || conf.organization === "") {
    conf.organization = DefaultOrganization
  }

  if (!conf.pullPolicy || conf.pullPolicy == "") {
    conf.pullPolicy = IfNotPresentPullPolicy
  }

  OpenAPI.BASE = `https://${conf.baseURL}/v1`

  if (conf.apiKey && conf.apiKey != "") {
    OpenAPI.HEADERS = {
      'Authorization': `Bearer ${conf.apiKey}`,
    }
  }

  return conf;
}

export function HexFromArrayBuffer (data: ArrayBuffer) {
  return Buffer.from(data).toString("hex");
}

export async function ComputeSHA256(data: ArrayBuffer): Promise<string> {
  return HexFromArrayBuffer(await crypto.subtle.digest("SHA-256", data))
}