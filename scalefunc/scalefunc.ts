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

// Version is the version of ScaleFunc
export type Version = string;

// V1Alpha is the V1 Alpha definition of a ScaleFunc
export const V1Alpha: Version = "v1alpha";

// V1Beta is the V1 Beta definition of a ScaleFunc
export const V1Beta: Version = "v1beta";

export type Language = string;

// Go is the Golang Source Language for Scale Functions
export const Go: Language = "go";

// Rust is the Rust Source Language for Scale Functions
export const Rust: Language = "rust";

// Typescript is the Typescript Source Language for Scale Functions
export const Typescript: Language = "ts";

// Javascript is the Javascript Source Language for Scale Functions
export const Javascript: Language = "js";

// AcceptedLanguages is an array of acceptable Languages
export const AcceptedLanguages: Language[] = [Go, Rust, Typescript, Javascript];

export function ValidString(str: string): boolean {
    return !InvalidStringRegex.test(str);
}

export class V1AlphaDependency {
    public Name: string;
    public Version: string;
    public Metadata: Map<string, string> | undefined;

    constructor(name: string, version: string, metadata: Map<string, string>) {
        this.Name = name;
        this.Version = version;
        this.Metadata = metadata;
    }
}

export class V1AlphaSchema {
    public Name: string;
    public Tag: string;
    public SignatureName: string;
    public SignatureBytes: Buffer;
    public SignatureHash: string;
    public Language: Language;
    public Stateless: boolean;
    public Dependencies: V1AlphaDependency[];
    public Function: Buffer;
    public Size: undefined | number;
    public Hash: undefined | string;

    constructor(
        name: string,
        tag: string,
        signatureName: string,
        signatureBytes: Buffer,
        signatureHash: string,
        language: Language,
        stateless: boolean,
        dependencies: V1AlphaDependency[],
        fn: Buffer
    ) {
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

    // Deprecated: Use V1BetaSchema instead
    Encode(): Uint8Array {
        const enc = new Encoder();
        enc.string(V1Alpha as string);

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

        const size = enc.bytes.length;
        const hashed = sha256(enc.bytes);
        const hex = Buffer.from(hashed).toString("hex");

        enc.uint32(size);
        enc.string(hex);

        return enc.bytes;
    }

    static Decode(data: Uint8Array): V1AlphaSchema {
        const dec = new Decoder(data);

        const version = dec.string() as Version;
        if (version !== V1Alpha) throw VersionErr;

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
        const dependencies: V1AlphaDependency[] = [];
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
            dependencies.push(new V1AlphaDependency(name, version, metadata));
        }

        const fn = dec.uint8Array();
        const size = dec.uint32();
        const hash = dec.string();

        const hashed = sha256(data.slice(0, size));
        const hex = Buffer.from(hashed).toString("hex");

        if (hex !== hash) throw ChecksumErr;

        const sf = new V1AlphaSchema(
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

// Deprecated: Use V1BetaSchema instead
export function ReadV1Alpha(path: string): V1AlphaSchema {
    return V1AlphaSchema.Decode(fs.readFileSync(path, null));
}

// Deprecated: Use V1BetaSchema instead
export function WriteV1Alpha(path: string, scaleFunc: V1AlphaSchema) {
    fs.writeFileSync(path, scaleFunc.Encode());
}

export class V1BetaExtension {
    public Name: string;
    public Organization: string;
    public Tag: string;
    public Schema: Buffer;
    public Hash: string;

    constructor(name: string, organization: string, tag: string, schema: Buffer, hash: string) {
        this.Name = name;
        this.Organization = organization;
        this.Tag = tag;
        this.Schema = schema;
        this.Hash = hash;
    }
}

export class V1BetaSignature {
    public Name: string;
    public Organization: string;
    public Tag: string;
    public Schema: Buffer;
    public Hash: string;

    constructor(name: string, organization: string, tag: string, schema: Buffer, hash: string) {
        this.Name = name;
        this.Organization = organization;
        this.Tag = tag;
        this.Schema = schema;
        this.Hash = hash;
    }
}

export class V1BetaSchema {
    public Name: string;
    public Tag: string;
    public Signature: V1BetaSignature;
    public Extensions: V1BetaExtension[];
    public Language: Language;
    public Manifest: Buffer;
    public Stateless: boolean;
    public Function: Buffer;
    public Size: undefined | number;
    public Hash: undefined | string;

    constructor(
        name: string,
        tag: string,
        signature: V1BetaSignature,
        extensions: V1BetaExtension[],
        language: Language,
        manifest: Buffer,
        stateless: boolean,
        fn: Buffer
    ) {
        this.Name = name;
        this.Tag = tag;
        this.Signature = signature;
        this.Extensions = extensions;
        this.Language = language;
        this.Manifest = manifest;
        this.Stateless = stateless;
        this.Function = fn;
    }

    Encode(): Uint8Array {
        const enc = new Encoder();
        enc.string(V1Beta as string);

        enc.string(this.Name);
        enc.string(this.Tag);

        enc.string(this.Signature.Name);
        enc.string(this.Signature.Organization);
        enc.string(this.Signature.Tag);
        enc.uint8Array(this.Signature.Schema);
        enc.string(this.Signature.Hash);

        enc.array(this.Extensions.length, Kind.Any);
        for (const ext of this.Extensions) {
            enc.string(ext.Name);
            enc.string(ext.Organization);
            enc.string(ext.Tag);
            enc.uint8Array(ext.Schema);
            enc.string(ext.Hash);
        }

        enc.string(this.Language as string);
        enc.uint8Array(this.Manifest);
        enc.boolean(this.Stateless);

        enc.uint8Array(this.Function);

        const size = enc.bytes.length;
        const hashed = sha256(enc.bytes);
        const hex = Buffer.from(hashed).toString("hex");

        enc.uint32(size);
        enc.string(hex);

        return enc.bytes;
    }

    static Decode(data: Uint8Array): V1BetaSchema {
        const dec = new Decoder(data);

        const version = dec.string() as Version;
        switch (version) {
            case V1Alpha:
               const v1Alpha = V1AlphaSchema.Decode(data);
                let orgSplit = v1Alpha.SignatureName.split("/");
                if (orgSplit.length === 1) {
                    orgSplit.unshift("");
                }

                let tagSplit = orgSplit[1].split(":");
                if (tagSplit.length === 1) {
                    tagSplit.push("");
                }

               const v1BetaSignature = new V1BetaSignature(tagSplit[0], orgSplit[0], tagSplit[1], v1Alpha.SignatureBytes, v1Alpha.SignatureHash);
               return new V1BetaSchema(v1Alpha.Name, v1Alpha.Tag, v1BetaSignature, [], v1Alpha.Language, Buffer.from([]), v1Alpha.Stateless, v1Alpha.Function);
            case V1Beta:
                break;
            default:
                throw VersionErr;
        }

        const name = dec.string();
        const tag = dec.string();
        const signatureName = dec.string();
        let signatureOrg = dec.string();
        let signatureTag = dec.string();
        const signatureBytes = dec.uint8Array();
        const signatureHash = dec.string();

        const v1BetaSignature = new V1BetaSignature(signatureName, signatureOrg, signatureTag, Buffer.from(signatureBytes), signatureHash);

        const extensionsSize = dec.array(Kind.Any);
        const extensions: V1BetaExtension[] = [];
        for (let i = 0; i < extensionsSize; i++) {
            const name = dec.string();
            const organization = dec.string();
            const tag = dec.string();
            const schema = dec.uint8Array();
            const hash = dec.string();
            extensions.push(new V1BetaExtension(name, organization, tag, Buffer.from(schema), hash));
        }

        const language = dec.string();
        if (!AcceptedLanguages.includes(language)) throw LanguageErr;

        const manifest = dec.uint8Array();

        const stateless = dec.boolean();

        const fn = dec.uint8Array();
        const size = dec.uint32();
        const hash = dec.string();

        const hashed = sha256(data.slice(0, size));
        const hex = Buffer.from(hashed).toString("hex");

        if (hex !== hash) throw ChecksumErr;

        const sf = new V1BetaSchema(
            name,
            tag,
            v1BetaSignature,
            extensions,
            language,
            Buffer.from(manifest),
            stateless,
            Buffer.from(fn)
        );
        sf.Size = size;
        sf.Hash = hash;
        return sf;
    }
}

export function ReadV1Beta(path: string): V1BetaSchema {
    return V1BetaSchema.Decode(fs.readFileSync(path, null));
}

export function WriteV1Beta(path: string, scaleFunc: V1BetaSchema) {
    fs.writeFileSync(path, scaleFunc.Encode());
}

export function Read(path: string): V1BetaSchema {
    return ReadV1Beta(path);
}

export function Write(path: string, scaleFunc: V1BetaSchema) {
    WriteV1Beta(path, scaleFunc);
}