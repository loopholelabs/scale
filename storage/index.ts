import { ScaleFunc, ValidString } from "@loopholelabs/scalefile/scalefunc"
import { Read, Write } from "@loopholelabs/scalefile/scalefunc/helpers"
import * as fs from "fs";
import * as os from "os";
import * as glob from "glob3";
import path from "path";

export const ErrInvalidName = new Error("invalid name");
export const ErrInvalidTag = new Error("invalid tag");
export const ErrInvalidOrg = new Error("invalid org");

export const DefaultCacheDirectory = ".config/scale/functions"

export interface Entry {
    scaleFunc: ScaleFunc,
    hash: string,
    organization: string,
}

export class Storage{
    public baseDirectory: string;
    constructor(baseDirectory: string) {
        this.baseDirectory = baseDirectory;
        if (!fs.existsSync(this.baseDirectory)) {
            fs.mkdirSync(this.baseDirectory, { recursive: true, mode: 0o755 });
        }
    }

    public Get(name: string, tag: string, org: string, hash: string | undefined): (Entry | undefined) {
        if (name.length === 0 || !ValidString(name)) {
            throw ErrInvalidName;
        }

        if (tag.length === 0 || !ValidString(tag)) {
            throw ErrInvalidTag;
        }

        if (org.length === 0 || !ValidString(org)) {
            throw ErrInvalidOrg;
        }

        if (hash) {
            const f = this.functionName(name, tag, org, hash);
            const p = this.fullPath(f);
            try {
                return {
                    hash: hash,
                    scaleFunc: Read(p),
                    organization: org,
                };
            } catch (error: any) {
                if (error.code === 'ENOENT') {
                    return undefined;
                }
                throw error
            }
        }

        const f = this.functionSearch(name, tag, org);
        const p = this.fullPath(f);
        const matches = glob.sync(p)
        if (matches.length === 0) {
            return undefined;
        }

        if (matches.length > 1) {
            throw new Error("multiple matches found for " + name + ":" + tag);
        }

        return {
            scaleFunc: Read(matches[0]),
            hash: path.basename(matches[0]).split("_")[3],
            organization: path.basename(matches[0]).split("_")[0],
        };
    }

    public Put(name: string, tag: string, org: string, hash: string, sf: ScaleFunc): void {
        const f = this.functionName(name, tag, org, hash);
        const p = this.fullPath(f);
        Write(p, sf);
    }

    public Delete(name: string, tag: string, org: string, hash: string): void {
        const f = this.functionName(name, tag, org, hash);
        const p = this.fullPath(f);
        fs.rmSync(p);
    }

    private fullPath(p: string): string {
        return path.join(this.baseDirectory, p);
    }

    private functionName(name: string, tag: string, org: string, hash: string): string {
        return `${org}_${name}_${tag}_${hash}_scale`;
    }

    private functionSearch(name: string, tag: string, org: string): string {
        return `${org}_${name}_${tag}_*_scale`;
    }
}

export const Default = new Storage(path.join(os.homedir(), DefaultCacheDirectory));