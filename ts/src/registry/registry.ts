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

import { ScaleFile } from "@loopholelabs/scalefile"
import { sha256 } from "js-sha256";

export type PullPolicy = string;

export const AlwaysPullPolicy: PullPolicy = "always";
export const IfNotPresentPullPolicy: PullPolicy = "if-not-present";
export const NeverPullPolicy: PullPolicy = "never";

export var ErrHasMismatch = new Error("hash mismatch");
export var ErrNoFunction = new Error("function does not exist and pull policy is never");

export interface Config {
  pullPolicy: PullPolicy
  cacheDirectory: string
  apiKey: string
}

export type Option = (c: Config) => void;

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

export function WithApiKey(apiKey: string): Option {
  return (c: Config) => {
    c.apiKey = apiKey;
  }
}

export async function NewPulldown(func: string, ...opts: [Option]): Promise<ScaleFile> {
  const config: Config = {
    pullPolicy: AlwaysPullPolicy,
    cacheDirectory: "~/.cache/scale/functions",
    apiKey: "",
  }

  for (const opt of opts) {
    opt(config);
  }

  let scaleFile = getFromCache(func, config);
  if (scaleFile && config.pullPolicy != AlwaysPullPolicy) {
    return scaleFile
  }

  // Contact the API endpoint with the request
  const response = apiRequest(func, config);
  const httpResp = await fetch(response.URL)
  const data = await httpResp.arrayBuffer();

  const hash = await computeSha256(data);
  if (hash !== response.Hash) {
    throw ErrHasMismatch;
  }


  const string = new TextDecoder("utf-8").decode(data);
  scaleFile = ScaleFile.Decode(string);
  saveToCache(func, scaleFile, config);
  return scaleFile
}

function buildFilename(func: string, config: Config): string {
    return `${func}.scale`;
}

function getFromCache (func: string, config: Config): ScaleFile {
  const file = buildFilename(func, config);
  const path = `${config.cacheDirectory}/${file}`;
  return ScaleFile.Read(path);
}

function saveToCache (func: string, scaleFile: ScaleFile, config: Config): void {
  const file = buildFilename(func, config);
  const path = `${config.cacheDirectory}/${file}`;
  ScaleFile.Write(path, scaleFile);
}

interface PulldownResponse {
  URL: string
  Hash: string
}

function apiRequest (func: string, config: Config): PulldownResponse {
  // TODO
  return {
    URL: "https://example.com",
    Hash: "deadbeef",
  }
}

async function computeSha256 (data: ArrayBuffer): Promise<string> {
  // use crypto subtle if available, otherwise fall back to a pure JS implementation
    if (window.crypto && window.crypto.subtle) {
      const hash = await crypto.subtle.digest("SHA-256", data);
      const hashArray = Array.from(new Uint8Array(hash));
      return hashArray.map((b) => b.toString(16).padStart(2, "0")).join("");
    }
    return sha256(data);
}
