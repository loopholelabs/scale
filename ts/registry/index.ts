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
import {sha256} from "js-sha256";
import {webcrypto as crypto} from "crypto";

import {Default, Storage} from "../storage";
import {models_GetFunctionResponse, OpenAPI, RegistryService} from "../client";
import https from "https";

export type PullPolicy = string;

export const AlwaysPullPolicy: PullPolicy = "always";
export const IfNotPresentPullPolicy: PullPolicy = "if-not-present";
export const NeverPullPolicy: PullPolicy = "never";

export var ErrHasMismatch = new Error("hash mismatch");
export var ErrNoFunction = new Error("function does not exist locally and pull policy does not allow pulling from registry");
export var ErrDownloadFailed = new Error("scale function could not be pull from the registry")

export const defaultBaseURL = "api.scale.sh"
export const defaultOrganization = "scale"

interface config {
  pullPolicy?: PullPolicy
  baseURL?: string,
  organization?: string,
  cacheDirectory?: string,
  apiKey?: string,
}

export type Option = (c: config) => void;

export function WithPullPolicy(pullPolicy: PullPolicy): Option {
  return (c: config) => {
    c.pullPolicy = pullPolicy;
  }
}

export function WithCacheDirectory(cacheDirectory: string): Option {
  return (c: config) => {
    c.cacheDirectory = cacheDirectory;
  }
}

export function WithAPIKey(apiKey: string): Option {
  return (c: config) => {
    c.apiKey = apiKey;
  }
}

export function WithBaseURL(baseURL: string): Option {
    return (c: config) => {
        c.baseURL = baseURL;
    }
}

export function WithOrganization(organization: string): Option {
    return (c: config) => {
        c.organization = organization;
    }
}

export async function New(name: string, tag: string, ...opts: Option[]): Promise<ScaleFunc> {
  const conf: config = {
    pullPolicy: IfNotPresentPullPolicy,
    baseURL: defaultBaseURL,
    organization: defaultOrganization,
  }

  for (const opt of opts) {
    opt(conf);
  }

  if (!conf.pullPolicy) {
    conf.pullPolicy = IfNotPresentPullPolicy
  }

  if (!conf.baseURL || conf.baseURL === "") {
    conf.baseURL = defaultBaseURL
  }

  if (!conf.organization || conf.organization === "") {
    conf.organization = defaultOrganization
  }

  let st = Default;
  if (conf.cacheDirectory && conf.cacheDirectory != "") {
    st = new Storage(conf.cacheDirectory)
  }

  OpenAPI.BASE = `https://${conf.baseURL}/v1`

  if (conf.apiKey && conf.apiKey != "") {
    OpenAPI.HEADERS = {
      'Authorization': `Bearer ${conf.apiKey}`,
    }
  }

  switch (conf.pullPolicy) {
    case NeverPullPolicy:
      const neverGet = st.Get(name, tag, conf.organization, "");
      if (!neverGet) {
        throw ErrNoFunction
      }
      return neverGet.scaleFunc;
    case IfNotPresentPullPolicy:
      const notPresentGet = st.Get(name, tag, conf.organization, "");
      if (notPresentGet) {
        return notPresentGet.scaleFunc
      }

      let notPresentFn: models_GetFunctionResponse | undefined = undefined;
      if (conf.organization === defaultOrganization) {
        notPresentFn = await RegistryService.getRegistryFunction1(name, tag)
      } else {
        notPresentFn = await RegistryService.getRegistryFunction3(name, tag, conf.organization)
      }

      if(!notPresentFn.presigned_url || !notPresentFn.hash) {
        throw ErrDownloadFailed
      }

      const notPresentData = await download(notPresentFn.presigned_url);
      const notPresentHash = await computeSHA256(notPresentData);
      if (notPresentHash !== notPresentFn.hash) {
        throw ErrHasMismatch;
      }
      const notPresentSF = ScaleFunc.Decode(new Uint8Array(notPresentData));
      st.Put(name, tag, conf.organization, notPresentHash, notPresentSF)
      return notPresentSF;
    case AlwaysPullPolicy:
      const alwaysGet = st.Get(name, tag, conf.organization, "");
      let alwaysFn: models_GetFunctionResponse | undefined = undefined;
      if (conf.organization === defaultOrganization) {
        alwaysFn = await RegistryService.getRegistryFunction1(name, tag)
      } else {
        alwaysFn = await RegistryService.getRegistryFunction3(name, tag, conf.organization)
      }

      if(!alwaysFn.presigned_url || !alwaysFn.hash) {
        throw ErrDownloadFailed
      }

      if (alwaysGet) {
        if (alwaysFn.hash === alwaysGet.hash) {
          return alwaysGet.scaleFunc
        }
      }

      const alwaysData = await download(alwaysFn.presigned_url);
      const alwaysHash = await computeSHA256(alwaysData);
      if (alwaysHash !== alwaysFn.hash) {
        throw ErrHasMismatch;
      }
      const alwaysSF = ScaleFunc.Decode(new Uint8Array(alwaysData));
      if (alwaysGet) {
        st.Delete(name, tag, alwaysGet.organization, alwaysGet.hash)
      }
      st.Put(name, tag, conf.organization, alwaysHash, alwaysSF)
      return alwaysSF;
    default:
      throw Error(`unknown pull policy ${conf.pullPolicy}`)
  }
}

async function download (url: string): Promise<ArrayBuffer> {
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

export async function computeSHA256 (data: ArrayBuffer): Promise<string> {
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
