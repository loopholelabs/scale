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

import { WasiContext } from "./runtime";

var barebonesWASI = function() {

  var WASI_ESUCCESS = 0;
  var WASI_EBADF = 8;
  var WASI_EINVAL = 28;
  var WASI_ENOSYS = 52;

  function fd_prestat_get(fd: number, bufPtr: number) {
      return WASI_EBADF;
  }

  function fd_prestat_dir_name(fd: number, pathPtr: number, pathLen: number) {
       return WASI_EINVAL;
  }

  function environ_sizes_get(environCount: number, environBufSize: number) {
      return WASI_ESUCCESS;
  }

  function environ_get(environ: number, environBuf: number) {
      return WASI_ESUCCESS;
  }

  function args_sizes_get(argc: number, argvBufSize: number) {
      return WASI_ESUCCESS;
  }

   function args_get(argv: number, argvBuf: number) {
      return WASI_ESUCCESS;
  }

  function fd_filestat_get() {
    return WASI_ENOSYS;
  }

  function fd_fdstat_get(fd: number, bufPtr: number) {
      return WASI_ESUCCESS;
  }

  function fd_read() {
    return WASI_ESUCCESS;
  }

  function fd_write(fd: number, iovs: number, iovsLen: number, nwritten: number) {
      return WASI_ESUCCESS;
  }

  function poll_oneoff(sin: number, sout: number, nsubscriptions: number, nevents: number) {
      return WASI_ENOSYS;
  }

  function proc_exit(rval: number) {
      return WASI_ENOSYS;
  }

  function fd_close(fd: number) {
      return WASI_ENOSYS;
  }

  function fd_seek(fd: number, offset: number, whence: number, newOffsetPtr: number) {
  }

  function path_open() {
    return WASI_ESUCCESS;
  }

  return {
      environ_sizes_get : environ_sizes_get,
      args_sizes_get : args_sizes_get,
      fd_prestat_get : fd_prestat_get,
      fd_fdstat_get : fd_fdstat_get,
      fd_filestat_get: fd_filestat_get,
      fd_read : fd_read,
      fd_write : fd_write,
      fd_prestat_dir_name : fd_prestat_dir_name,
      path_open: path_open,
      environ_get : environ_get,
      args_get : args_get,
      poll_oneoff : poll_oneoff,
      proc_exit : proc_exit,
      fd_close : fd_close,
      fd_seek : fd_seek,
      clock_time_get : () => WASI_ESUCCESS
  }               
}


export function getNewWasi(): WasiContext {
  const w: WasiContext = {
    getImportObject: () =>{  /* Minimal import object */
      return barebonesWASI();
      /*{
        fd_write: () => {},
        args_sizes_get: () => 0,
        args_get: () => {},
        clock_time_get: () => {},
        environ_get: () => 0,
        environ_sizes_get: () => 0,
        fd_close: () => 0,
        fd_fdstat_get: () => 0,
        fd_seek: () => 0,
        proc_exit: () => 0,
      };*/
    },
    start: (instance: WebAssembly.Instance) => {
      const startFn = (instance.exports._start as Function | undefined);
      if (startFn) {
        startFn();
      }
    }
  }
  return w;
}