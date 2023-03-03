# Typescript guest builder

## Quickstart

1. Build the rust jsbuilder

```  
  cd builder
  ./install-wasi-sdk.sh
  make
```

2. Build the example http scale functions

```    
  cd runner
  npm install
  ./build.sh
```

This should create 3 wasm modules.

| Scale file              | Description                              |
| ----------------------- | ---------------------------------------- |
| module_error.wasm       | This scale function just throws an error |
| module_endpoint.wasm    | This scale function sets the Body |
| module_middleware.wasm  | This scale function adds a header |
