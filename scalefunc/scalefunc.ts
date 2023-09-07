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

import sha256 from "fast-sha256";

import {Decoder, Encoder, Kind} from "@loopholelabs/polyglot";
import fs from "fs";

// eslint-disable-next-line @typescript-eslint/no-var-requires
const Buffer = require("buffer/").Buffer;

export const InvalidStringRegex = /[^A-Za-z0-9-.]/;

export const VersionErr = new Error("unknown or invalid version");
export const LanguageErr = new Error("unknown or invalid language");
export const ChecksumErr = new Error("error while verifying checksum");

export type Version = string;

// V1Alpha is the V1 Alpha definition of a ScaleFunc
export const V1Alpha: Version = "v1alpha";

export type Language = string;

// Go is the Golang Source Language for Scale Functions
export const Go: Language = "go";

// Rust is the Rust Source Language for Scale Functions
export const Rust: Language = "rust";

// Typescript is the Typescript Source Language for Scale Functions
export const Typescript: Language = "ts";

// Javascript is the Javascript Source Language for Scale Functions
export const Javascript: Language = "js";

// AcceptedVersions is an array of acceptable Versions
export const AcceptedVersions: Version[] = [V1Alpha];

// AcceptedLanguages is an array of acceptable Languages
export const AcceptedLanguages: Language[] = [Go, Rust, Typescript, Javascript];

export class Dependency {
    public Name: string;
    public Version: string;
    public Metadata: Map<string, string> | undefined;

    constructor(name: string, version: string, metadata: Map<string, string>) {
        this.Name = name;
        this.Version = version;
        this.Metadata = metadata;
    }
}

export class ScaleFunc {
    public Version: Version;
    public Name: string;
    public Tag: string;
    public SignatureName: string;
    public SignatureBytes: Buffer;
    public SignatureHash: string;
    public Language: Language;
    public Stateless: boolean;
    public Dependencies: Dependency[];
    public Function: Buffer;
    public Size: undefined | number;
    public Hash: undefined | string;

    constructor(
        version: Version,
        name: string,
        tag: string,
        signatureName: string,
        signatureBytes: Buffer,
        signatureHash: string,
        language: Language,
        stateless: boolean,
        dependencies: Dependency[],
        fn: Buffer
    ) {
        this.Version = version;
        this.Name = name;
        this.Tag = tag;
        this.SignatureName = signatureName;
        this.SignatureBytes = signatureBytes;
        this.SignatureHash = signatureHash;
        this.Language = language;
        this.Stateless = stateless;
        this.Dependencies = dependencies;
        this.Function = fn;
    }

    Encode(): Uint8Array {
        const enc = new Encoder();
        enc.string(this.Version as string);
        enc.string(this.Name);
        enc.string(this.Tag);
        enc.string(this.SignatureName);
        enc.uint8Array(this.SignatureBytes);
        enc.string(this.SignatureHash);
        enc.string(this.Language as string);
        enc.boolean(this.Stateless);

        enc.array(this.Dependencies.length, Kind.Any);
        for (const d of this.Dependencies) {
            enc.string(d.Name);
            enc.string(d.Version);
            enc.map(d.Metadata?.size ?? 0, Kind.String, Kind.String);
            for (const [k, v] of d.Metadata ?? []) {
                enc.string(k);
                enc.string(v);
            }
        }

        enc.uint8Array(this.Function);

        // Compute the hash (sha256)

        const size = enc.bytes.length;
        const hashed = sha256(enc.bytes);
        const hex = Buffer.from(hashed).toString("hex");

        enc.uint32(size);
        enc.string(hex);

        return enc.bytes;
    }

    static Decode(data: Uint8Array): ScaleFunc {
        const dec = new Decoder(data);

        const version = dec.string() as Version;
        if (!AcceptedVersions.includes(version)) throw VersionErr;

        const name = dec.string();
        const tag = dec.string();
        const signatureName = dec.string();
        const signatureBytes = dec.uint8Array();
        const signatureHash = dec.string();

        const language = dec.string();
        if (!AcceptedLanguages.includes(language)) throw LanguageErr;

        let stateless = false;
        try {
            stateless = dec.boolean();
        } catch (_) {} // eslint-disable-line no-empty

        const dependenciesSize = dec.array(Kind.Any);
        const dependencies: Dependency[] = [];
        for (let i = 0; i < dependenciesSize; i++) {
            const name = dec.string();
            const version = dec.string();
            const metadataSize = dec.map(Kind.String, Kind.String);
            const metadata = new Map<string, string>();
            for (let j = 0; j < metadataSize; j++) {
                const key = dec.string();
                const value = dec.string();
                metadata.set(key, value);
            }
            dependencies.push(new Dependency(name, version, metadata));
        }

        const fn = dec.uint8Array();
        const size = dec.uint32();
        const hash = dec.string();

        const hashed = sha256(data.slice(0, size));
        const hex = Buffer.from(hashed).toString("hex");

        if (hex !== hash) throw ChecksumErr;

        const sf = new ScaleFunc(
            version,
            name,
            tag,
            signatureName,
            Buffer.from(signatureBytes),
            signatureHash,
            language,
            stateless,
            dependencies,
            Buffer.from(fn)
        );
        sf.Size = size;
        sf.Hash = hash;
        return sf;
    }
}

export function ValidString(str: string): boolean {
    return !InvalidStringRegex.test(str);
}

export function Read(path: string): ScaleFunc {
    return ScaleFunc.Decode(fs.readFileSync(path, null));
}

export function Write(path: string, scaleFunc: ScaleFunc) {
    fs.writeFileSync(path, scaleFunc.Encode());
}