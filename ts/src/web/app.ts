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
import { ModuleBrowser } from '../runtime/moduleBrowser';
import { Host } from '../runtime/host';
import { Context, Request, Response, StringList } from "../runtime/generated/generated";
import { Context as ourContext} from '../runtime/context';


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

const runButton = (document.getElementById('crun') as HTMLInputElement);

runButton.onclick = async function() {

    console.log("Creating a context");

    // Create a context to send in...
    // TODO: Read from form

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

    Host.showContext(context);

    // TODO: Read from form...

    const modHttpEndpoint = await fetch('./http-endpoint.wasm');
    const arrayHttpEndpoint = await modHttpEndpoint.arrayBuffer();
    let moduleHttpEndpoint = new ModuleBrowser(Buffer.from(arrayHttpEndpoint), null);

    const modHttpMiddleware = await fetch('./http-middleware.wasm');
    const arrayHttpMiddleware = await modHttpMiddleware.arrayBuffer();
    let moduleHttpMiddleware = new ModuleBrowser(Buffer.from(arrayHttpMiddleware), moduleHttpEndpoint);

    let ctx = new ourContext(context);

    let ctime = (new Date()).getTime();
    let retContext = moduleHttpMiddleware.run(ctx);
    let etime = (new Date()).getTime();
    
    Host.showContext(retContext.context());
    // TODO: Show in the fields etc...

    (document.getElementById('status') as HTMLInputElement).innerHTML = "Executed in " + (etime - ctime).toFixed(3) + "ms"

    let resp = retContext.context().Response;

    (document.getElementById('rstatus') as HTMLInputElement).value = "" + resp.StatusCode;
    let dec = new TextDecoder();
    (document.getElementById('rbody') as HTMLInputElement).value = dec.decode(resp.Body);

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

