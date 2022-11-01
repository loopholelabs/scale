package io.loopholelabs.scale;

import org.teavm.interop.Import;
import org.teavm.interop.Export;

public class Client {

	@Import(name = "next", module = "env")
	public static native int nextModule(int ptr, int size);

  @Export(name = "malloc")
  public static int malloc(int size) {
    return 1;
  }

  @Export(name = "run")
  public static int run(int ptr, int size) {
    // TODO: Read the context from memory, and decode... Do something to it, encode, store it in memory.

    //    System.out.println("Run was called with " + ptr + ", " + size);

    // nextModule()
    return 1234;
  }

	// Compile without, and this class won't be included. Hm.
	static void main(String... args) {}
}
