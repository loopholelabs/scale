# Typescript scale host

The typescript scale host can be used either as a node library, or in a browser context.

## Run tests

`npm run test`

## Build

`npm run build`

## Browser

`npm run start`

Browse to `http://localhost:1234` and play with scale functions.

## Adapters

### Express
The express adapter `ExpressAdapter` can be used as middleware.

    app.use(new ExpressAdapter(runtime));

### http / fetch API
There are also conversion utils for http and fetchAPI included.

## Quickstart using module

Embedding the runtime is very simple. For a full example see `runtime.test.ts` or the `web/app.ts`.

    // Load a wasm module up
    const modWasm = fs.readFileSync("./example_modules/http-endpoint.wasm");

    // Create a Module and init()
    const moduleHttpEndpoint = new Module(modWasm, getNewWasi());
    await moduleHttpEndpoint.init();

    // Create a runtime with a list of scale functions
    const runtime = new Runtime([moduleHttpMiddleware]);

    // Create a context, and run the runtime

    const ctx = new Context(context);
    const retContext = runtime.run(ctx);
    // retContext is the output from the scale function chain.
