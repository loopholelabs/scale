import * as fs from 'fs';
import { Module } from './runtime/module';
import { Context, Request, Response, StringList } from "./runtime/generated/generated";
import { Host } from './runtime/host';
import { Context as ourContext} from './runtime/context';

// Create a context to send in...
var enc = new TextEncoder();
let body = enc.encode("Hello world this is a request body");
let headers = new Map<string, StringList>();
headers.set('content', new StringList(['hello']));
let req1 = new Request('GET', BigInt(100), 'https', '1.2.3.4', body, headers);
let respBody = enc.encode("Response body");
let respHeaders = new Map<string, StringList>();        
const resp1 = new Response(200, respBody, respHeaders);        
const context = new Context(req1, resp1);

// Now we can use context...

const modHttpEndpoint = fs.readFileSync('./example_modules/http-endpoint.wasm');
const modHttpMiddleware = fs.readFileSync('./example_modules/http-middleware.wasm');
let moduleHttpEndpoint = new Module(modHttpEndpoint, null);
let moduleHttpMiddleware = new Module(modHttpMiddleware, moduleHttpEndpoint);

// Run the modules...

let ctx = new ourContext(context);

console.log("\nINPUT CONTEXT")
Host.showContext(context);

let retContext = moduleHttpMiddleware.run(ctx);

console.log("\nOUTPUT CONTEXT");
Host.showContext(retContext.context());
