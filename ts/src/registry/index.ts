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

import {ScaleFunc} from "@loopholelabs/scalefile"
import { sha256 } from "js-sha256";
import * as fs from "fs";
import { webcrypto as crypto } from "crypto";
import https from "https";
import * as os from "os";
import path from "path";

export type PullPolicy = string;

export const AlwaysPullPolicy: PullPolicy = "always";
export const IfNotPresentPullPolicy: PullPolicy = "if-not-present";
export const NeverPullPolicy: PullPolicy = "never";

export var ErrHasMismatch = new Error("hash mismatch");
export var ErrNoFunction = new Error("function does not exist and pull policy is never");
export var ErrDownloadFailed = new Error("The scalefunc could not be retrieved from the server")

const defaultApiBaseUrl = "https://api.scale.sh/v1"
const defaultCacheSubpath = ".cache/scale/functions"

export interface Config {
  pullPolicy: PullPolicy
  cacheDirectory: string
  apiKey: string,
  apiBaseUrl: string,
  organization: string,
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

export function WithApiBaseUrl(apiBaseUrl: string): Option {
    return (c: Config) => {
        c.apiBaseUrl = apiBaseUrl;
    }
}

export function WithOrganization(organization: string): Option {
    return (c: Config) => {
        c.organization = organization;
    }
}

export async function New(func: string, tag: string, ...opts: Option[]): Promise<ScaleFunc> {
  const defaultCacheDirectory = path.join(os.homedir(), defaultCacheSubpath)

  const config: Config = {
    pullPolicy: AlwaysPullPolicy,
    cacheDirectory: defaultCacheDirectory,
    apiKey: "",
    apiBaseUrl: defaultApiBaseUrl,
    organization: "default",
  }

  for (const opt of opts) {
    opt(config);
  }

  let scaleFunc: ScaleFunc | undefined = undefined;

  if (config.pullPolicy !== AlwaysPullPolicy) {
    let scaleFunc = getFromCache(func, tag, undefined, config);
    if (scaleFunc && config.pullPolicy != AlwaysPullPolicy) {
      return scaleFunc
    }

    if (config.pullPolicy === NeverPullPolicy) {
      throw ErrNoFunction;
    }
  }

  // Contact the API endpoint with the request
  const response = await apiRequest(func, tag, config);
  scaleFunc = getFromCache(func, tag, response.hash, config);
  if (scaleFunc) {
    return scaleFunc;
  }

  const data = await downloadScaleFunc(response.presigned_url);

  const hash = await computeSha256(data);
  if (hash !== response.hash) {
    throw ErrHasMismatch;
  }

  scaleFunc = ScaleFunc.Decode(new Uint8Array(data));
  saveToCache(func, tag, hash, config, scaleFunc);
  return scaleFunc;
}

function buildFilename(func: string, tag: string, hash: string, config: Config): string {
    return `${func}.${tag}.${hash}.scale`;
}

function getFromCache (func: string, tag: string, hash: string | undefined, config: Config): ScaleFunc | undefined {
  if (hash) {
    const file = buildFilename(func, tag, hash, config);
    const filePath = path.join(config.cacheDirectory, file);
    try {
      return ScaleFunc.Read(filePath);
    } catch (error: any) {
      if (error.code === 'ENOENT') {
        return undefined;
      }
      throw error
    }
  }

  const filePrefix = `${func}.${tag}.`
    const files = fs.readdirSync(config.cacheDirectory);
    for (const file of files) {
      if (file.startsWith(filePrefix)) {
        const filePath = path.join(config.cacheDirectory, file);
        try {
          return ScaleFunc.Read(filePath);
        } catch (error: any) {
          if (error.code === 'ENOENT') {
            return undefined;
          }
          throw error
        }
      }
    }
}

export function saveToCache (func: string, tag: string, hash: string, config: Config, scaleFunc: ScaleFunc): void {
  const file = buildFilename(func, tag, hash, config);
  if (!fs.existsSync(config.cacheDirectory)) {
    fs.mkdirSync(config.cacheDirectory, { recursive: true });
  }
  const path = `${config.cacheDirectory}/${file}`;
  ScaleFunc.Write(path, scaleFunc);
}

interface GetFunctionResponse {
  name: string,
  tag: string,
  organization: string,
  public: boolean,
  hash: string,
  presigned_url: string,
}

async function apiRequest (func: string, tag: string, config: Config): Promise<GetFunctionResponse> {
  if (typeof window !== "undefined") {
    const response = await fetch(`${config.apiBaseUrl}/registry/function/${config.organization}/${func}/${tag}`, {
      method: 'get',
      headers: new Headers({
        'Authorization': `Bearer ${config.apiKey}`,
      }),
    });
    return response.json();
  }
  return new Promise((resolve, reject) => {
    let body: Uint8Array[] = [];
    const request = https.request(`${config.apiBaseUrl}/registry/function/${config.organization}/${func}/${tag}`, {
      method: 'get',
      headers: {
        'Authorization': `Bearer ${config.apiKey}`,
      },
    }, (response) => {
      response.on('data', (chunk) => {
        body.push(chunk);
      })

      response.on('end', () => {
        const data = Buffer.concat(body).toString();
        resolve(JSON.parse(data));
      })

      response.on('error', (error) => {
        reject(error);
      });
    });

    request.end();
  });
}

async function downloadScaleFunc (url: string): Promise<ArrayBuffer> {
  if (typeof window !== "undefined") {
    const httpResp = await fetch(url)
    return httpResp.arrayBuffer();
  }
  return new Promise((resolve, reject) => {
    let body: Uint8Array[] = [];
    const request = https.request(url, (response) => {
      if (response.statusCode !== 200) {
        reject(ErrDownloadFailed);
      }

      response.on('data', (chunk) => {
        body.push(chunk);
      })

      response.on('end', () => {
        resolve(Buffer.concat(body));
      })

      response.on('error', (error) => {
        reject(error);
      });
    });

    request.end();
  });
}

export async function computeSha256 (data: ArrayBuffer): Promise<string> {
  // use crypto subtle if available, otherwise fall back to a pure JS implementation
  if (crypto && crypto.subtle) {
    const hash = await crypto.subtle.digest("SHA-256", data);
    return base64FromArrayBuffer(hash);
  }
  const hash = sha256.create();
  hash.update(data);
  return base64FromArrayBuffer(hash.arrayBuffer())
}

/*
* This method, while appearing quite unorthodox, is the most efficient way to get a base64 value from an array buffer
* while avoiding an expensive intermediate string representation
* */
async function base64FromArrayBuffer (data: ArrayBuffer) {
  if (typeof window !== "undefined") {
    const base64url: string = await new Promise((r) => {
      const reader = new FileReader()
      reader.onload = () => r(reader.result as string)
      reader.readAsDataURL(new Blob([data]))
    })

    return base64url.split(",", 2)[1]
  } else {
    return Buffer.from(data).toString("base64").replace(/\+/g, "-");
  }
}
