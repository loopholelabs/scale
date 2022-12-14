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

//import { TextEncoder, TextDecoder } from 'util';
//import * as fs from 'fs';
import { Context, Request, Response, StringList } from "../http-signature/generated/generated";
import { Runtime, WasiContext } from '../sigruntime/runtime';

import { ScaleFunc, V1Alpha, Go } from "@loopholelabs/scalefile";

import { HttpContext, HttpContextFactory } from "../http-signature/HttpContext";

const addHeaderButton = (document.getElementById('cheadersadd') as HTMLInputElement);
addHeaderButton.onclick = function() {
    // Add a new row...
    const cheaders = (document.getElementById('cheaders') as HTMLTableElement);

    const newRow = cheaders.insertRow(0);
    const cell1 = newRow.insertCell(0);
    const cell2 = newRow.insertCell(1);
    const cell3 = newRow.insertCell(2);

    const inputKey = document.createElement("input");
    inputKey.type = "text";
    inputKey.value = "KEY";
    cell1.appendChild(inputKey);

    const inputVal = document.createElement("input");
    inputVal.type = "text";
    inputVal.value = "VALUE";
    cell2.appendChild(inputVal);
    
    const deleteButton = document.createElement("button");
    deleteButton.appendChild(document.createTextNode("Delete"));
    deleteButton.onclick = function(r) {
        return function() {
            r.remove();
        }
    }(newRow);

    cell3.appendChild(deleteButton);
}


const addButton = (document.getElementById('cadd') as HTMLInputElement);

addButton.onclick = async function() {
  let inputfile = (document.getElementById('inputfile') as HTMLInputElement);
  if (inputfile.files!=null) {
    let file = inputfile.files[0];
    console.log(file);

    let reader = new FileReader();

    reader.readAsArrayBuffer(file);
  
    reader.onload = function() {
      console.log(reader.result);
      const d = Buffer.from(reader.result as ArrayBuffer);
      const scalefn = new ScaleFunc(V1Alpha, file.name, "ExampleName@ExampleVersion", Go, [], d);
  
      addModule(scalefn);
    };
  
    reader.onerror = function() {
      console.log(reader.error);
    };
  }
}

var barebonesWASI = function() {

  var moduleInstanceExports: any = null;

  var WASI_ESUCCESS = 0;
  var WASI_EBADF = 8;
  var WASI_EINVAL = 28;
  var WASI_ENOSYS = 52;

  var WASI_STDOUT_FILENO = 1;

  function setModuleInstance(instance: any) {

      moduleInstanceExports = instance.exports;
  }

  function getModuleMemoryDataView() {
      // call this any time you'll be reading or writing to a module's memory 
      // the returned DataView tends to be dissaociated with the module's memory buffer at the will of the WebAssembly engine 
      // cache the returned DataView at your own peril!!

      return new DataView(moduleInstanceExports.memory.buffer);
  }

  function fd_prestat_get(fd: number, bufPtr: number) {
      return WASI_EBADF;
  }

  function fd_prestat_dir_name(fd: number, pathPtr: number, pathLen: number) {
       return WASI_EINVAL;
  }

  function environ_sizes_get(environCount: number, environBufSize: number) {
/*
      var view = getModuleMemoryDataView();

      view.setUint32(environCount, 0, !0);
      view.setUint32(environBufSize, 0, !0);
*/
      return WASI_ESUCCESS;
  }

  function environ_get(environ: number, environBuf: number) {

      return WASI_ESUCCESS;
  }

  function args_sizes_get(argc: number, argvBufSize: number) {
/*
      var view = getModuleMemoryDataView();

      view.setUint32(argc, 0, !0);
      view.setUint32(argvBufSize, 0, !0);
*/
      return WASI_ESUCCESS;
  }

   function args_get(argv: number, argvBuf: number) {

      return WASI_ESUCCESS;
  }

  function fd_fdstat_get(fd: number, bufPtr: number) {
/*
      var view = getModuleMemoryDataView();

      view.setUint8(bufPtr, fd);
      view.setUint16(bufPtr + 2, 0, !0);
      view.setUint16(bufPtr + 4, 0, !0);

      function setBigUint64(byteOffset: number, value: number, littleEndian: boolean) {

          var lowWord = value;
          var highWord = 0;

          view.setUint32(littleEndian ? 0 : 4, lowWord, littleEndian);
          view.setUint32(littleEndian ? 4 : 0, highWord, littleEndian);
     }

      setBigUint64(bufPtr + 8, 0, !0);
      setBigUint64(bufPtr + 8 + 8, 0, !0);
*/
      return WASI_ESUCCESS;
  }

  function fd_write(fd: number, iovs: number, iovsLen: number, nwritten: number) {
/*
      var view = getModuleMemoryDataView();

      var written = 0;
      var bufferBytes: any[] = [];                   

      function getiovs(iovs: number, iovsLen: number) {
          // iovs* -> [iov, iov, ...]
          // __wasi_ciovec_t {
          //   void* buf,
          //   size_t buf_len,
          // }
          var buffers = Array.from({ length: iovsLen }, function (_, i) {
                 var ptr = iovs + i * 8;
                 var buf = view.getUint32(ptr, !0);
                 var bufLen = view.getUint32(ptr + 4, !0);

                 return new Uint8Array(moduleInstanceExports.memory.buffer, buf, bufLen);
              });

          return buffers;
      }

      var buffers = getiovs(iovs, iovsLen);
      function writev(iov: Uint8Array) {

          for (var b = 0; b < iov.byteLength; b++) {

             bufferBytes.push(iov[b]);
          }

          written += b;
      }

      buffers.forEach(writev);

      if (fd === WASI_STDOUT_FILENO) console.log(String.fromCharCode.apply(null, bufferBytes));                            

      view.setUint32(nwritten, written, !0);
*/
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

  return {
      setModuleInstance : setModuleInstance,
      environ_sizes_get : environ_sizes_get,
      args_sizes_get : args_sizes_get,
      fd_prestat_get : fd_prestat_get,
      fd_fdstat_get : fd_fdstat_get,
      fd_write : fd_write,
      fd_prestat_dir_name : fd_prestat_dir_name,
      environ_get : environ_get,
      args_get : args_get,
      poll_oneoff : poll_oneoff,
      proc_exit : proc_exit,
      fd_close : fd_close,
      fd_seek : fd_seek,
      clock_time_get : () => WASI_ESUCCESS
  }               
}


function getNewWasi() {

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
      console.log("EXPORTS", instance.exports);

      const startFn = (instance.exports._start as Function | undefined);
      if (startFn) {
        console.log("Call _start on wasm module...");
        startFn();
      }
    }
  }
  return w;

}

const runButton = (document.getElementById('crun') as HTMLInputElement);

let modules: ScaleFunc[] = [];

async function init() {
  const examples = [
    "./middleware-java.wasm",
    "./middleware-typescript.wasm",
    "./go-middleware.wasm",
    "./go-endpoint.wasm",
    "./java-endpoint.wasm"
  ];

  for(let i=0;i<examples.length;i++) {
    const mod = await fetch(examples[i]);
    const arrayModule = await mod.arrayBuffer();
    const scalefn = new ScaleFunc(V1Alpha, examples[i], "ExampleName@ExampleVersion", Go, [], Buffer.from(arrayModule));
    addModule(scalefn);
  }
}

init();

function addModule(m: ScaleFunc) {
  const tab = (document.getElementById("cmodules") as HTMLTableElement);

  const newRow = tab.tBodies[0].insertRow(-1);
  const cell1 = newRow.insertCell(0);
  cell1.appendChild(document.createTextNode(m.Name===undefined?"":m.Name));

  const cell2 = newRow.insertCell(1);
  const delbutton = document.createElement("a");
  delbutton.href = "#";
  delbutton.className = "delbutton";
  delbutton.appendChild(document.createTextNode("Delete"));
  cell2.appendChild(delbutton);


  delbutton.onclick = function(mod, row) {
    return function() {
      row.remove();
      // Delete this module from the array, and from the UI...
      const index = modules.indexOf(mod);
      if (index > -1) {
        modules.splice(index, 1);
        console.log("Removed module from array ", modules);
      }
    }
  }(m, newRow);

  modules.push(m);
}


runButton.onclick = async function() {

    console.log("Creating a context");

    // Create a context to send in...

    const method = (document.getElementById('cmethod') as HTMLInputElement).value;
    const protocol = (document.getElementById('cprotocol') as HTMLInputElement).value;
    const ip = (document.getElementById('cip') as HTMLInputElement).value;
    const body = (document.getElementById('cbody') as HTMLInputElement).value;

    let enc = new TextEncoder();
    let bodyData = enc.encode(body);
    let headers = new Map<string, StringList>();

    const cheaders = (document.getElementById('cheaders') as HTMLTableElement);

    let heads: Map<string, string[]> = new Map();

    for (let i=0;i<cheaders.rows.length;i++) {
        let row = cheaders.rows[i];
        let ikey = (row.cells[0].firstChild as HTMLInputElement).value;
        let ival = (row.cells[1].firstChild as HTMLInputElement).value;
        console.log("TODO: " + ikey + " = " + ival);
        if (heads.has(ikey)) {
            heads.get(ikey)?.push(ival);
        } else {
            heads.set(ikey, [ival]);
        }
    }

    for (let k of heads.keys()) {
        let vals = heads.get(k);
        if (vals===undefined) {
            vals = [];
        }
        headers.set(k, new StringList(vals));
    }
    let req1 = new Request(method, BigInt(bodyData.length), protocol, ip, bodyData, headers);
    
    let respBody = enc.encode("Response body");
    let respHeaders = new Map<string, StringList>();        
    const resp1 = new Response(200, respBody, respHeaders);        
    const context = new Context(req1, resp1);

    const signatureFactory = HttpContextFactory;

    const r = new Runtime<HttpContext>(getNewWasi, signatureFactory, modules);
    await r.Ready;

    const i = await r.Instance(null);
    i.Context().ctx = context;

    let ctime = (new Date()).getTime();
    i.Run();
    let etime = (new Date()).getTime();

    const retContext = i.Context().ctx;

    if (retContext!=null) {

      (document.getElementById('status') as HTMLInputElement).innerHTML = "Executed in " + (etime - ctime).toFixed(3) + "ms"

      let resp = retContext.Response;

      (document.getElementById('rstatus') as HTMLElement).innerHTML = "" + resp.StatusCode;
      let dec = new TextDecoder();
      let rbody = (document.getElementById('rbody') as HTMLElement);
      while(rbody.childNodes.length>0) rbody.removeChild(rbody.childNodes[0]);
      rbody.appendChild(document.createTextNode(dec.decode(resp.Body)));

      const rheaders = (document.getElementById('rheaders') as HTMLTableElement);

      while(rheaders.rows.length>0) {
          rheaders.deleteRow(0);
      }

      for (let k of resp.Headers.keys()) {
          const newRow = rheaders.insertRow(0);
          const cell1 = newRow.insertCell(0);
          cell1.appendChild(document.createTextNode(k));
          const cell2 = newRow.insertCell(1);

          let values = resp.Headers.get(k);
          if (values!=undefined) {
              for (let i of values.Value.values()) {
                  cell2.appendChild(document.createTextNode(i));
                  cell2.appendChild(document.createElement("br"));
              }
          }
      }
    }
}
