# Adapters

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

