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

import { Writer } from "./config";

const ErrNoInstance = Error("no webassembly instance")
const ErrNoMemory = Error("no exported memory in webassembly instance")

export interface RequiredFunctions extends WebAssembly.ModuleImports {
    fd_prestat_get(fd: number, bufPtr: number): number
    fd_prestat_dir_name(fd: number, pathPtr: number, pathLen: number): number
    environ_sizes_get(environCount: number, environBufSize: number): number
    environ_get(environ: number, environBuf: number): number
    args_sizes_get(argc: number, argvBufSize: number): number
    args_get(argv: number, argvBuf: number): number
    fd_filestat_get(fd: number, bufPtr: number): number
    fd_fdstat_get(fd: number, bufPtr: number): number
    fd_read(fd: number, iovs_ptr: number, iovs_len: number, nread_ptr: number): number
    fd_write(fd: number, iovs: number, iovsLen: number, nwritten: number): number
    fd_close(fd: number): number
    fd_seek(fd: number, offset: number, whence: number, newOffsetPtr: number): number
    poll_oneoff(sin: number, sout: number, nsubscriptions: number, nevents: number): number
    proc_exit(rval: number): number
    path_open(fd: number, dirflags: number, path_ptr: number, path_len: number, oflags: number, fs_rights_base: number, fs_rights_inheriting: number, fd_flags: number, opened_fd_ptr: number): number
    clock_time_get(id: number, precision: bigint, time: number): number
    random_get(buf: number, buf_len: number): void
}

export class DisabledWASI {
    public static ESUCCESS = 0;
    public static EBADF = 8;
    public static EINVAL = 28;
    public static ENOSYS = 52;

    public static CLOCKID_REALTIME = 0;

    private exports: WebAssembly.Exports | undefined;

    private readonly env: { [key: string]: string } | undefined;
    private readonly stdout: Writer | undefined;
    private readonly stderr: Writer | undefined;

    constructor(env?: { [key: string]: string }, stdout?: Writer, stderr?: Writer) {
        this.env = env;
        this.stdout = stdout;
        this.stderr = stderr;
    }

    private getDataView(): DataView {
        if (!this.exports) {
            throw ErrNoInstance;
        }
        if (!this.exports.memory) {
            throw ErrNoMemory;
        }

        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-ignore
        return new DataView(this.exports.memory.buffer);
    }

    public SetInstance(instance: WebAssembly.Instance) {
        this.exports = instance.exports;
    }

    public environ_sizes_get(environCount: number, environBufSize: number): number {
        const view = this.getDataView();
        let count = 0;
        let size = 0;
        if (this.env !== undefined) {
            for (const key in this.env) {
                count++;
                size += key.length + 1 + this.env[key].length + 1;
            }
        }
        view.setUint32(environCount, count, true);
        view.setUint32(environBufSize, size, true);
        return DisabledWASI.ESUCCESS;
    }

    public environ_get(environ: number, environBuf: number): number {
        const view = this.getDataView();
        let offset = 0;
        if (this.env !== undefined) {
            for (const key in this.env) {
                const value = this.env[key];
                const keyLen = key.length;
                const valueLen = value.length;
                view.setUint32(environ + offset, environBuf, true);
                offset += 4;
                view.setUint32(environ + offset, keyLen + 1 + valueLen + 1, true);
                offset += 4;
                for (let i = 0; i < keyLen; i++) {
                    view.setUint8(environBuf, key.charCodeAt(i));
                    environBuf++;
                }
                view.setUint8(environBuf, 0);
                environBuf++;
                for (let i = 0; i < valueLen; i++) {
                    view.setUint8(environBuf, value.charCodeAt(i));
                    environBuf++;
                }
                view.setUint8(environBuf, 0);
                environBuf++;
            }
        }
        return DisabledWASI.ESUCCESS;
    }

    public args_sizes_get(argc: number, argvBufSize: number): number {
        const view = this.getDataView();
        view.setUint32(argc, 0, true);
        view.setUint32(argvBufSize, 0, true);
        return DisabledWASI.ESUCCESS;
    }

    public args_get(_argv: number, _argvBuf: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.ESUCCESS;
    }

    public fd_prestat_get(_fd: number, _bufPtr: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.EBADF;
    }

    public fd_prestat_dir_name(_fd: number, _pathPtr: number, _pathLen: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.EINVAL;
    }

    public fd_filestat_get(_fd: number, _bufPtr: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.EBADF;
    }

    public fd_fdstat_get(_fd: number, _bufPtr: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.EBADF;
    }

    public fd_read(_fd: number, _iovs_ptr: number, _iovs_len: number, _nread_ptr: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.EBADF;
    }

    public fd_write(fd: number, iovsPtr: number, iovsLength: number, bytesWrittenPtr: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        if ( (fd === 1 && this.stdout !== undefined) || (fd === 2 && this.stderr !== undefined) ) {
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            const iovs = new Uint32Array(this.exports.memory.buffer, iovsPtr, iovsLength * 2);
            let text = "";
            let totalBytesWritten = 0;
            const decoder = new TextDecoder();
            for(let i =0; i < iovsLength * 2; i += 2){
                const offset = iovs[i];
                const length = iovs[i+1];
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-ignore
                const textChunk = decoder.decode(new Int8Array(this.exports.memory.buffer, offset, length));
                text += textChunk;
                totalBytesWritten += length;
            }
            const dataView = this.getDataView();
            dataView.setInt32(bytesWrittenPtr, totalBytesWritten, true);
            if (fd === 1 && this.stdout !== undefined) {
                this.stdout(text);
            } else if (fd === 2 && this.stderr !== undefined) {
                this.stderr(text);
            }

            return DisabledWASI.ESUCCESS;
        }

        return DisabledWASI.EBADF;
    }

    public fd_close(_fd: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.EBADF;
    }

    public fd_seek(_fd: number, _offset: number, _whence: number, _newOffsetPtr: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.EBADF;
    }

    public poll_oneoff(_sin: number, _sout: number, _nsubscriptions: number, _nevents: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.ENOSYS;
    }

    public proc_exit(_rval: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.ENOSYS;
    }

    public path_open(_fd: number, _dirflags: number, _path_ptr: number, _path_len: number, _oflags: number, _fs_rights_base: number, _fs_rights_inheriting: number, _fd_flags: number, _opened_fd_ptr: number): number { // eslint-disable-line @typescript-eslint/no-unused-vars
        return DisabledWASI.EBADF;
    }

    public clock_time_get(id: number, _precision: bigint, time: number) { // eslint-disable-line @typescript-eslint/no-unused-vars
        const buffer = this.getDataView()
        if (id === DisabledWASI.CLOCKID_REALTIME) {
            buffer.setBigUint64(time, BigInt(new Date().getTime()) * 1000000n, true);
        } else {
            buffer.setBigUint64(time, 0n, true);
        }
        return DisabledWASI.ESUCCESS;
    }

    public random_get(buf: number, buf_len: number): void {
        const buffer = this.getDataView();
        for (let i = 0; i < buf_len; i++) {
            buffer.setInt8(buf+i, (Math.random() * 256) | 0);
        }
    }

    public GetImports(): RequiredFunctions {
        return {
            environ_sizes_get: this.environ_sizes_get.bind(this),
            environ_get: this.environ_get.bind(this),
            args_sizes_get: this.args_sizes_get.bind(this),
            args_get: this.args_get.bind(this),
            fd_prestat_get: this.fd_prestat_get.bind(this),
            fd_prestat_dir_name: this.fd_prestat_dir_name.bind(this),
            fd_filestat_get: this.fd_filestat_get.bind(this),
            fd_fdstat_get: this.fd_fdstat_get.bind(this),
            fd_read: this.fd_read.bind(this),
            fd_write: this.fd_write.bind(this),
            fd_close: this.fd_close.bind(this),
            fd_seek: this.fd_seek.bind(this),
            poll_oneoff: this.poll_oneoff.bind(this),
            proc_exit: this.proc_exit.bind(this),
            path_open: this.path_open.bind(this),
            clock_time_get: this.clock_time_get.bind(this),
            random_get: this.random_get.bind(this),
        }
    }
}
