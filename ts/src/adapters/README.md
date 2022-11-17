# Adapters

First, you must create a `Runtime` encapsulating the Scale Functions you wish to run.

    const modHttpEndpoint = fs.readFileSync(
      "./example_modules/http-endpoint.wasm"
    );
    const modHttpMiddleware = fs.readFileSync(
      "./example_modules/http-middleware.wasm"
    );

    const moduleHttpEndpoint = new Module(modHttpEndpoint, getNewWasi());
    await moduleHttpEndpoint.init();

    const moduleHttpMiddleware = new Module(modHttpMiddleware, getNewWasi());
    await moduleHttpMiddleware.init();

    const runtime = new Runtime([moduleHttpMiddleware, moduleHttpEndpoint]);


## ExpressAdapter

The express adapter is simple to use. It can be used as middleware

    const app = express();
    const adapter = new ExpressAdapter(runtime);
    app.use(adapter.getHandler());

## HttpAdapter

If you are using http you can do a similar setup.

    var adapter = new HttpAdapter(runtime);
    var server = http.createServer(adapter.getHandler());

## NextAdapter

There is also an adapter for next.js edge functions as per https://vercel.com/docs/concepts/functions/edge-functions#creating-edge-functions

    var nextAdapter = new NextAdapter(runtime);

    export default nextAdapter.getHandler();

    export const config = {
      runtime: 'experimental-edge',
    };
