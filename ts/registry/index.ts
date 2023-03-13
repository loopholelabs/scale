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

import {ScaleFunc} from "@loopholelabs/scalefile/scalefunc"
import {sha256} from "js-sha256";
import {webcrypto as crypto} from "crypto";

import {Default, Storage} from "../storage";
import {models_GetFunctionResponse, OpenAPI, RegistryService} from "../client";
import https from "https";

import {
  ComputeSHA256 as BrowserComputeSHA256,
  DefaultConfig,
  Option,
  HexFromArrayBuffer,
  ErrHasMismatch,
  ErrDownloadFailed,
  ErrNoFunction,
  DefaultBaseURL,
  DefaultOrganization,
  PullPolicy,
  NeverPullPolicy,
  IfNotPresentPullPolicy,
  AlwaysPullPolicy,
  WithAPIKey,
  WithBaseURL,
  WithOrganization,
  WithPullPolicy,
  WithCacheDirectory,
} from "./browser"

export { ErrHasMismatch, ErrDownloadFailed, ErrNoFunction }
export { PullPolicy, NeverPullPolicy, IfNotPresentPullPolicy, AlwaysPullPolicy }
export { DefaultBaseURL, DefaultOrganization }

export { WithAPIKey, WithBaseURL, WithOrganization, WithCacheDirectory, WithPullPolicy }

export async function Download(name: string, tag: string, ...opts: Option[]): Promise<ScaleFunc> {
  const conf = DefaultConfig(...opts)

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
      if (conf.organization === DefaultOrganization) {
        notPresentFn = await RegistryService.getRegistryFunction1(name, tag)
      } else {
        notPresentFn = await RegistryService.getRegistryFunction3(conf.organization, name, tag)
      }

      if(!notPresentFn.presigned_url || !notPresentFn.hash) {
        throw ErrDownloadFailed
      }

      const notPresentData = await RawDownload(notPresentFn.presigned_url);
      const notPresentHash = await ComputeSHA256(notPresentData);
      if (notPresentHash !== notPresentFn.hash) {
        throw ErrHasMismatch;
      }
      const notPresentSF = ScaleFunc.Decode(new Uint8Array(notPresentData));
      st.Put(name, tag, conf.organization, notPresentHash, notPresentSF)
      return notPresentSF;
    case AlwaysPullPolicy:
      const alwaysGet = st.Get(name, tag, conf.organization, "");
      let alwaysFn: models_GetFunctionResponse | undefined = undefined;
      if (conf.organization === DefaultOrganization) {
        alwaysFn = await RegistryService.getRegistryFunction1(name, tag)
      } else {
        alwaysFn = await RegistryService.getRegistryFunction3(conf.organization, name, tag)
      }

      if(!alwaysFn.presigned_url || !alwaysFn.hash) {
        throw ErrDownloadFailed
      }

      if (alwaysGet) {
        if (alwaysFn.hash === alwaysGet.hash) {
          return alwaysGet.scaleFunc
        }
      }

      const alwaysData = await RawDownload(alwaysFn.presigned_url);
      const alwaysHash = await ComputeSHA256(alwaysData);
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

export async function ComputeSHA256 (data: ArrayBuffer): Promise<string> {
  if (crypto && crypto.subtle) {
    return BrowserComputeSHA256(data);
  }
  const hash = sha256.create();
  hash.update(data);
  return HexFromArrayBuffer(hash.arrayBuffer())
}

async function RawDownload (url: string): Promise<ArrayBuffer> {
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
