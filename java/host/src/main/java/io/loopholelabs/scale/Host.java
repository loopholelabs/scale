package io.loopholelabs.scale;

import java.io.*;
import java.nio.file.*;

import org.wasmer.*;

public class Host {

  public static void main(String[] args) throws Exception {
    // `simple.wasm` is located at `tests/resources/`.
    Path wasmPath = Paths.get("example_modules/java-endpoint.wasm");

    // Reads the WebAssembly module as bytes.
    byte[] wasmBytes = Files.readAllBytes(wasmPath);

    System.out.println("Loaded some wasm " + wasmBytes.length);

    // Instantiates the WebAssembly module.
    Instance instance = new Instance(wasmBytes);

    System.out.println("Wasm instance " + instance);

    System.out.println("Exports" + instance.exports);

    // Calls an exported function, and returns an object array.
    Object[] results = instance.exports.getFunction("sum").apply(5, 37);

    System.out.println((Integer) results[0]); // 42

    // Drops an instance object pointer which is stored in Rust.
    instance.close();
  }

}