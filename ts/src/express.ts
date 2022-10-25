//import * as express from 'express';

import express = require('express');
import bodyParser = require('body-parser');

import * as fs from 'fs';
import { Module } from './runtime/module';
import { Context, Request, Response, StringList } from "./runtime/generated/generated";
import { Host } from './runtime/host';
import { Context as ourContext} from './runtime/context';
import { ExpressAdapter } from './adapters/expressAdapter';

var app = express();

let port = 8090;

const modHttpEndpoint = fs.readFileSync('./example_modules/http-endpoint.wasm');
const modHttpMiddleware = fs.readFileSync('./example_modules/http-middleware.wasm');
let moduleHttpEndpoint = new Module(modHttpEndpoint, null);
let moduleHttpMiddleware = new Module(modHttpMiddleware, moduleHttpEndpoint);

var adapter = new ExpressAdapter(moduleHttpMiddleware);

app.use(bodyParser.raw({
    type: (t)=>true,
}));
app.use(adapter.handler.bind(adapter));
  
app.listen(port, () => {
    console.log(`Example app listening on port ${port}`)
})