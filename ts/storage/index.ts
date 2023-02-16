import {ScaleFunc} from "@loopholelabs/scalefile"
import * as fs from "fs";
import * as os from "os";
import * as glob from "glob";

// @ts-ignore
import path from "path";

const defaultCacheDirectory = ".config/scale/functions"

export class Storage{
    public baseDirectory: string;
    constructor(baseDirectory: string) {
        this.baseDirectory = baseDirectory;
        if (!fs.existsSync(this.baseDirectory)) {
            fs.mkdirSync(this.baseDirectory, { recursive: true, mode: 0o755 });
        }
    }

    public Get(name: string, tag: string, org: string, hash: string | undefined): ({
        hash: string,
        scaleFunc: ScaleFunc
    } | undefined) {
        if (hash) {
            const f = this.functionName(name, tag, org, hash);
            const p = this.fullPath(f);
            try {
                return {
                    hash: hash,
                    scaleFunc: ScaleFunc.Read(p)
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
            hash: matches[0].split(".")[3],
            scaleFunc: ScaleFunc.Read(matches[0])
        };
    }

    public Put(name: string, tag: string, org: string, hash: string, sf: ScaleFunc): void {
        const f = this.functionName(name, tag, org, hash);
        const p = this.fullPath(f);
        ScaleFunc.Write(p, sf);
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
        return `${org}.${name}.${tag}.${hash}.scale`;
    }

    private functionSearch(name: string, tag: string, org: string): string {
        return `${org}.${name}.${tag}.*.scale`;
    }
}

export const Default = new Storage(path.join(os.homedir(), defaultCacheDirectory));