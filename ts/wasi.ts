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

import { TextDecoder } from "util";

const WASI_ESUCCESS = 0;
const WASI_EBADF = 8;
const WASI_EINVAL = 28;
const WASI_ENOSYS = 52;

class disabledWASI {
  private instance: WebAssembly.Instance | undefined;

  public setInstance(i: WebAssembly.Instance) {
    this.instance = i;
  }

  fd_write(fd: number, iovs: number, iovsLen: number, nwritten: number) {
    if (this.instance == undefined) {
      return WASI_EBADF;
    }
    // Look at the memory, and update the number written...
    const mem = this.instance.exports.memory as WebAssembly.Memory;
    const memData = new Uint8Array(mem.buffer);

    // Get the io vectors
    let bytesWritten = 0;
    let iovs_ptr = iovs;

    for (let vec = 0;vec < iovsLen; vec ++) {
      const v = memData.slice(iovs_ptr, iovs_ptr + 8);
      const dv = new DataView(v.buffer);
      const ptr = dv.getUint32(0, true);
      const len = dv.getUint32(4, true);
      iovs_ptr += 8;
      bytesWritten += len;

      const b = memData.slice(ptr, ptr + len);
      const dec = new TextDecoder();
      //const hex = [...b].map(x => x.toString(16).padStart(2, '0')).join('');

      // NOTE: The console.log/error functions add a newline at the end.
      if (fd == 1) {
        console.log(dec.decode(b));
      } else if (fd == 2) {
        console.error(dec.decode(b));
      }
    }

    new DataView(mem.buffer).setUint32(nwritten, bytesWritten, true);
    return WASI_ESUCCESS;
  }

  fd_prestat_get(fd: number, bufPtr: number) {
    return WASI_EBADF;
  }

  fd_prestat_dir_name(fd: number, pathPtr: number, pathLen: number) {
    return WASI_EINVAL;
  }

  environ_sizes_get(environCount: number, environBufSize: number) {
    return WASI_ESUCCESS;
  }

  environ_get(environ: number, environBuf: number) {
    return WASI_ESUCCESS;
  }

  args_sizes_get(argc: number, argvBufSize: number) {
    return WASI_ESUCCESS;
  }

  args_get(argv: number, argvBuf: number) {
    return WASI_ESUCCESS;
  }

  fd_filestat_get() {
    return WASI_ENOSYS;
  }

  fd_fdstat_get(fd: number, bufPtr: number) {
    return WASI_ESUCCESS;
  }

  fd_read() {
    return WASI_ESUCCESS;
  }

  poll_oneoff(sin: number, sout: number, nsubscriptions: number, nevents: number) {
    return WASI_ENOSYS;
  }

  proc_exit(rval: number) {
    return WASI_ENOSYS;
  }

  fd_close(fd: number) {
    return WASI_ENOSYS;
  }

  fd_seek(fd: number, offset: number, whence: number, newOffsetPtr: number) {
  }

  path_open() {
    return WASI_ESUCCESS;
  }

  clock_time_get() {
    return WASI_ESUCCESS;
  }

  public imports(): any {
    return {
        environ_sizes_get: this.environ_sizes_get.bind(this),
        args_sizes_get: this.args_sizes_get.bind(this),
        fd_prestat_get: this.fd_prestat_get.bind(this),
        fd_fdstat_get: this.fd_fdstat_get.bind(this),
        fd_filestat_get: this.fd_filestat_get.bind(this),
        fd_read: this.fd_read.bind(this),
        fd_write: this.fd_write.bind(this),
        fd_prestat_dir_name: this.fd_prestat_dir_name.bind(this),
        path_open: this.path_open.bind(this),
        environ_get: this.environ_get.bind(this),
        args_get: this.args_get.bind(this),
        poll_oneoff: this.poll_oneoff.bind(this),
        proc_exit: this.proc_exit.bind(this),
        fd_close: this.fd_close.bind(this),
        fd_seek: this.fd_seek.bind(this),
        clock_time_get: this.clock_time_get.bind(this),
    }
  }
}


interface WASIContext {
    start(instance: WebAssembly.Instance): void;
    getImportObject(): any;
}

export function NewWASI(): WASIContext {
    let w = new disabledWASI();

    return {
      getImportObject: () => {
          return w.imports();
      },
      start: (instance: WebAssembly.Instance) => {
          w.setInstance(instance);

          const startFn = (instance.exports._start as Function | undefined);
          if (startFn) {
              startFn();
          }
      }
  };
}
